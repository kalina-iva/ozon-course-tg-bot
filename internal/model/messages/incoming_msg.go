package messages

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/helper/tracelog"
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

type ExpenseRepository interface {
	Create(ctx context.Context, userID int64, category string, amount uint64, date time.Time) error
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

type reportProducer interface {
	SendReportMessage(userID int64, msg []byte) error
}

type messageSender interface {
	SendMessage(text string, cases []string, userID int64) error
}

type reportRequest struct {
	UserID       int64     `json:"user_id"`
	Period       time.Time `json:"period"`
	CurrencyCode string    `json:"currency_code"`
}

type Model struct {
	tgClient         messageSender
	expenseRepo      ExpenseRepository
	exchangeRateRepo exchangeRateRepository
	userRepo         userRepository
	txManager        txManager
	reportProducer   reportProducer
}

func New(
	tgClient messageSender,
	expenseRepo ExpenseRepository,
	exchangeRateRepo exchangeRateRepository,
	userRepo userRepository,
	txManager txManager,
	producer reportProducer,
) *Model {
	return &Model{
		tgClient:         tgClient,
		expenseRepo:      expenseRepo,
		exchangeRateRepo: exchangeRateRepo,
		userRepo:         userRepo,
		txManager:        txManager,
		reportProducer:   producer,
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

func (m *Model) IncomingMessage(ctx context.Context, msg Message) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "start processing message")
	defer span.Finish()
	tracelog.Info(span, "start tracing incoming message")

	var text string
	var cases []string

	params := strings.Split(msg.Text, " ")
	switch params[0] {
	case "/start":
		text = manual
	case "/help":
		text = manual
	case "/newexpense":
		text = m.newExpenseHandler(ctx, msg.UserID, params)
	case "/report":
		text = m.reportHandler(ctx, msg.UserID, params)
	case "/setcurrency":
		text = chooseCurrency
		cases = AvailableCurrencies
	case "/setlimit":
		text = m.limitHandler(ctx, msg.UserID, params)
	case "/dellimit":
		text = m.delLimitHandler(ctx, msg.UserID)
	default:
		logger.Info("unknown command", zap.String("command", params[0]))
		text = unknownCommand
	}
	return params[0], m.tgClient.SendMessage(text, cases, msg.UserID)
}

func (m *Model) SetCurrency(ctx context.Context, msg CallbackQuery) error {
	err := m.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
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

func (m *Model) newExpenseHandler(ctx context.Context, userID int64, params []string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "expense handler")
	defer span.Finish()

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

	err = m.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
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
		if amount, err = m.convertAmountToRub(ctx, *user, parsedAmount); err != nil {
			return err
		}
		if err = m.expenseRepo.Create(ctx, userID, category, amount, date); err != nil {
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

func (m *Model) convertAmountToRub(ctx context.Context, user entity.User, amount float64) (uint64, error) {
	code := m.getCurrencyCode(user)
	rate, err := m.exchangeRateRepo.GetRate(ctx, code)
	if err != nil {
		return 0, errors.Wrap(err, "get exchange rate")
	}
	amount /= rate
	return uint64(math.Round(amount * cntKopInRub)), nil
}

func (m *Model) SendReport(ctx context.Context, report string, userID int64) error {
	return m.tgClient.SendMessage(report, nil, userID)
}

func (m *Model) reportHandler(ctx context.Context, userID int64, params []string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "report handler")
	defer span.Finish()

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

	user, err := m.userRepo.GetUser(ctx, userID)
	if err != nil {
		logger.Error("cannot get user", zap.Error(err))
		return userNotFound
	}

	err = m.sendToQueue(user, period)
	if err != nil {
		logger.Error("cannot send message to report queue", zap.Error(err))
		return canNotCreateReport
	}

	return reportIsGenerated
}

func (m *Model) sendToQueue(user *entity.User, period time.Time) error {
	msg, err := json.Marshal(reportRequest{
		UserID:       user.ID,
		Period:       period,
		CurrencyCode: m.getCurrencyCode(*user),
	})
	if err != nil {
		return errors.Wrap(err, "cannot marshal report message")
	}

	err = m.reportProducer.SendReportMessage(user.ID, msg)
	if err != nil {
		return errors.Wrap(err, "cannot send report message to queue")
	}
	return nil
}

func (m *Model) getCurrencyCode(user entity.User) string {
	if user.CurrencyCode != nil {
		return *user.CurrencyCode
	}
	return DefaultCurrencyCode
}

func (m *Model) limitHandler(ctx context.Context, userID int64, params []string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "limit handler")
	defer span.Finish()

	parsedAmount, err := m.parseAmount(params[1])
	if err != nil {
		logger.Error("cannot parse amount", zap.Error(err))
		return invalidAmount
	}
	err = m.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := m.userRepo.GetUser(ctx, userID)
		if err != nil {
			logger.Error("cannot get user", zap.Error(err))
			return errUserNotFound
		}
		var amount uint64
		if amount, err = m.convertAmountToRub(ctx, *user, parsedAmount); err != nil {
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

func (m *Model) delLimitHandler(ctx context.Context, userID int64) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "del limit handler")
	defer span.Finish()

	err := m.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
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
