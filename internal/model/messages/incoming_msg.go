package messages

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

const (
	cntKopInRub         = 100
	DefaultCurrencyCode = "RUB"
)

var AvailableCurrencies = []string{"RUB", "USD", "EUR", "CNY"}

type exchangeRateRepository interface {
	GetRate(ctx context.Context, code string) (float64, error)
}

type expenseRepository interface {
	New(ctx context.Context, userID int64, category string, amount uint64, date time.Time) error
	Report(ctx context.Context, userID int64, period time.Time) ([]*entity.Report, error)
	GetAmountByPeriod(ctx context.Context, userID int64, period time.Time) (uint64, error)
}

type userRepository interface {
	GetUser(ctx context.Context, userID int64) (*entity.User, error)
	SetCurrency(ctx context.Context, userID int64, currency string) error
	SetLimit(ctx context.Context, userID int64, limit uint64) error
	DelLimit(ctx context.Context, userID int64) error
}

type txManager interface {
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

type messageSender interface {
	SendMessage(text string, cases []string, userID int64) error
}

type Model struct {
	tgClient         messageSender
	expenseRepo      expenseRepository
	exchangeRateRepo exchangeRateRepository
	userRepo         userRepository
	txManager        txManager
	ctx              context.Context
}

func New(
	ctx context.Context,
	tgClient messageSender,
	expenseRepo expenseRepository,
	exchangeRateRepo exchangeRateRepository,
	userRepo userRepository,
	txManager txManager,
) *Model {
	return &Model{
		ctx:              ctx,
		tgClient:         tgClient,
		expenseRepo:      expenseRepo,
		exchangeRateRepo: exchangeRateRepo,
		userRepo:         userRepo,
		txManager:        txManager,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type CallbackQuery struct {
	Data   string
	UserID int64
}

func (m *Model) IncomingMessage(msg Message) (string, error) {
	/*
		var span opentracing.Span
		span, m.ctx = opentracing.StartSpanFromContext(m.ctx, "operation_name")
		defer span.Finish()

		if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
			logger.Info("trace",
				zap.String("id", spanContext.TraceID().String()),
			)
		}
	*/

	var text string
	var cases []string

	params := strings.Split(msg.Text, " ")
	switch params[0] {
	case "/start":
		text = manual
	case "/help":
		text = manual
	case "/newexpense":
		text = m.newExpenseHandler(msg.UserID, params)
	case "/report":
		text = m.reportHandler(msg.UserID, params)
	case "/setcurrency":
		text = chooseCurrency
		cases = AvailableCurrencies
	case "/setlimit":
		text = m.limitHandler(msg.UserID, params)
	case "/dellimit":
		text = m.delLimitHandler(msg.UserID)
	default:
		logger.Info("unknown command", zap.String("command", params[0]))
		text = unknownCommand
	}
	return params[0], m.tgClient.SendMessage(text, cases, msg.UserID)
}

func (m *Model) SetCurrency(msg CallbackQuery) error {
	err := m.txManager.WithinTransaction(m.ctx, func(ctx context.Context) error {
		_, err := m.userRepo.GetUser(ctx, msg.UserID)
		if err != nil {
			logger.Error("cannot get user", zap.Error(err))
			return errUserNotFound
		}
		if err := m.userRepo.SetCurrency(ctx, msg.UserID, msg.Data); err != nil {
			return err
		}
		return nil
	})

	var text string
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			text = userNotFound
		} else {
			logger.Error("cannot set currency", zap.Error(err))
			text = canNotSaveCurrency
		}
	} else {
		text = currencySaved
	}

	return m.tgClient.SendMessage(text, nil, msg.UserID)
}

func (m *Model) newExpenseHandler(userID int64, params []string) string {
	const cntRequiredParams = 3
	if len(params) < cntRequiredParams {
		logger.Info("no required params", zap.Strings("params", params))
		return needCategoryAndAmount
	}
	category := params[1]

	parsedAmount, err := m.parseAmount(params[2])
	if err != nil {
		logger.Error("cannot parse amount", zap.Error(err))
		return invalidAmount
	}

	var date time.Time
	if len(params) == cntRequiredParams+1 {
		date, err = time.Parse("01-02-2006", params[3])
		if err != nil {
			logger.Error("cannot parse date", zap.Error(err))
			return invalidDate
		}
	} else {
		date = time.Now()
	}

	err = m.txManager.WithinTransaction(m.ctx, func(ctx context.Context) error {
		user, err := m.userRepo.GetUser(ctx, userID)
		if err != nil {
			logger.Error("cannot get user", zap.Error(err))
			return errUserNotFound
		}
		if user.MonthlyLimit != nil {
			var currentSum uint64
			if currentSum, err = m.expenseRepo.GetAmountByPeriod(ctx, userID, beginningOfMonth()); err != nil {
				return err
			}
			if currentSum > *user.MonthlyLimit {
				logger.Info("limit exceeded", zap.Uint64("current", currentSum), zap.Uint64("limit", *user.MonthlyLimit))
				return errLimitExceeded
			}
		}
		var amount uint64
		if amount, err = m.convertAmountToRub(*user, parsedAmount); err != nil {
			return err
		}
		if err = m.expenseRepo.New(ctx, userID, category, amount, date); err != nil {
			return err
		}
		return nil
	})

	return getMsgTextForExpenseByErr(err)
}

func getMsgTextForExpenseByErr(err error) string {
	if err == nil {
		return expenseAdded
	}
	if errors.Is(err, errLimitExceeded) {
		return limitExceeded
	}
	if errors.Is(err, errUserNotFound) {
		return userNotFound
	}
	logger.Error("cannot get sum for period", zap.Error(err))
	return canNotAddExpense
}

func beginningOfMonth() time.Time {
	now := time.Now()
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

func (m *Model) parseAmount(amountStr string) (float64, error) {
	const bitSize = 64
	amount, err := strconv.ParseFloat(amountStr, bitSize)
	if err != nil {
		return 0, errors.Wrap(err, "parse amount")
	}
	if amount <= 0 {
		return 0, errInvalidAmount
	}
	return amount, nil
}

func (m *Model) convertAmountToRub(user entity.User, amount float64) (uint64, error) {
	code := m.getCurrencyCode(user)
	rate, err := m.exchangeRateRepo.GetRate(m.ctx, code)
	if err != nil {
		return 0, errors.Wrap(err, "get exchange rate")
	}
	amount /= rate
	return uint64(math.Round(amount * cntKopInRub)), nil
}

func (m *Model) reportHandler(userID int64, params []string) string {
	const cntRequiredParams = 2
	if len(params) < cntRequiredParams {
		logger.Info("no required params", zap.Strings("params", params))
		return needPeriod
	}
	var period time.Time
	now := time.Now()
	switch params[1] {
	case "y":
		period = now.AddDate(-1, 0, 0)
	case "m":
		period = now.AddDate(0, -1, 0)
	case "w":
		period = now.AddDate(0, 0, -7)
	default:
		logger.Info("invalid period", zap.String("period", params[1]))
		return invalidPeriod
	}

	user, err := m.userRepo.GetUser(m.ctx, userID)
	if err != nil {
		logger.Error("cannot get user", zap.Error(err))
		return userNotFound
	}
	code := m.getCurrencyCode(*user)
	rate, err := m.exchangeRateRepo.GetRate(m.ctx, code)
	if err != nil {
		logger.Error("cannot get rate from expenseRepo", zap.Error(err))
		return canNotGetRate
	}

	var sb strings.Builder
	currencyShort := getCurrencyShortByCode(code)
	report, err := m.expenseRepo.Report(m.ctx, userID, period)
	if err != nil {
		logger.Error("cannot get report", zap.Error(err))
		return canNotCreateReport
	}

	for _, item := range report {
		amount := float64(item.AmountInKopecks) * rate
		sb.WriteString(item.Category)
		sb.WriteString(fmt.Sprintf(": %.2f %v\n", amount/cntKopInRub, currencyShort))
	}
	if sb.Len() == 0 {
		logger.Info("no data for report")
		return noDataForReport
	}
	return sb.String()
}

func (m *Model) getCurrencyCode(user entity.User) string {
	if user.CurrencyCode != nil {
		return *user.CurrencyCode
	}
	return DefaultCurrencyCode
}

func getCurrencyShortByCode(code string) (short string) {
	switch code {
	case "RUB":
		short = "₽"
	case "USD":
		short = "＄"
	case "EUR":
		short = "€"
	case "CNY":
		short = "元"
	}
	return
}

func (m *Model) limitHandler(userID int64, params []string) string {
	parsedAmount, err := m.parseAmount(params[1])
	if err != nil {
		logger.Error("cannot parse amount", zap.Error(err))
		return invalidAmount
	}
	err = m.txManager.WithinTransaction(m.ctx, func(ctx context.Context) error {
		user, err := m.userRepo.GetUser(ctx, userID)
		if err != nil {
			logger.Error("cannot get user", zap.Error(err))
			return errUserNotFound
		}
		var amount uint64
		if amount, err = m.convertAmountToRub(*user, parsedAmount); err != nil {
			return err
		}
		if err := m.userRepo.SetLimit(ctx, userID, amount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			return userNotFound
		}
		logger.Error("cannot save limit", zap.Error(err))
		return canNotSaveLimit
	}
	return limitSaved
}

func (m *Model) delLimitHandler(userID int64) string {
	err := m.txManager.WithinTransaction(m.ctx, func(ctx context.Context) error {
		_, err := m.userRepo.GetUser(ctx, userID)
		if err != nil {
			logger.Error("cannot get user", zap.Error(err))
			return errUserNotFound
		}
		if err := m.userRepo.DelLimit(ctx, userID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			return userNotFound
		}
		logger.Error("cannot save limit", zap.Error(err))
		return canNotSaveLimit
	}
	return limitDeleted
}
