package messages

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

const (
	cntKopInRub         = 100
	DefaultCurrencyCode = "RUB"
)

var AvailableCurrencies = []string{"RUB", "USD", "EUR", "CNY"}

type exchangeRateRepository interface {
	GetRate(code string) (float64, error)
}

type expenseRepository interface {
	New(ctx context.Context, userID int64, category string, amount uint64, date time.Time) error
	Report(ctx context.Context, userID int64, period time.Time) ([]*entity.Report, error)
	GetAmountByPeriod(ctx context.Context, userID int64, period time.Time) (uint64, error)
}

type userRepository interface {
	SetCurrency(userID int64, currency string) error
	GetCurrency(userID int64) *string
	SetLimit(userID int64, limit uint64) error
	DelLimit(userID int64) error
	GetLimit(userID int64) *uint64
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

func (m *Model) IncomingMessage(msg Message) (err error) {
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
		text = unknownCommand
	}
	return m.tgClient.SendMessage(text, cases, msg.UserID)
}

func (m *Model) SetCurrency(msg CallbackQuery) error {
	if err := m.userRepo.SetCurrency(msg.UserID, msg.Data); err != nil {
		log.Println("cannot set currency:", err)
		return m.tgClient.SendMessage(canNotSaveCurrency, nil, msg.UserID)
	}
	return m.tgClient.SendMessage(currencySaved, nil, msg.UserID)
}

func (m *Model) newExpenseHandler(userID int64, params []string) string {
	const cntRequiredParams = 3
	if len(params) < cntRequiredParams {
		return needCategoryAndAmount
	}
	category := params[1]
	amount, err := m.parseAmount(userID, params[2])
	if err != nil {
		log.Println("cannot parse amount:", err)
		return invalidAmount
	}

	var date time.Time
	if len(params) == cntRequiredParams+1 {
		date, err = time.Parse("01-02-2006", params[3])
		if err != nil {
			log.Println("cannot parse date:", err)
			return invalidDate
		}
	} else {
		date = time.Now()
	}

	err = m.txManager.WithinTransaction(m.ctx, func(ctx context.Context) error {
		limit := m.userRepo.GetLimit(userID)
		if limit != nil {
			currentSum, err := m.expenseRepo.GetAmountByPeriod(ctx, userID, beginningOfMonth())
			if err != nil {
				return err
			}
			if currentSum > *limit {
				return errLimitExceeded
			}
		}
		if err = m.expenseRepo.New(m.ctx, userID, category, amount, date); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.As(err, &errLimitExceeded) {
			return limitExceeded
		}
		log.Println("cannot get sum for period:", err)
		return canNotAddExpense
	}

	return expenseAdded
}

func beginningOfMonth() time.Time {
	now := time.Now()
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

func (m *Model) parseAmount(userID int64, amountStr string) (uint64, error) {
	const bitSize = 64
	amount, err := strconv.ParseFloat(amountStr, bitSize)
	if err != nil {
		return 0, errors.Wrap(err, "parse amount")
	}
	if amount <= 0 {
		return 0, errInvalidAmount
	}
	code := m.getCurrencyCode(userID)
	rate, err := m.exchangeRateRepo.GetRate(code)
	if err != nil {
		return 0, errors.Wrap(err, "get exchange rate")
	}
	amount /= rate
	return uint64(math.Round(amount * cntKopInRub)), nil
}

func (m *Model) reportHandler(userID int64, params []string) string {
	const cntRequiredParams = 2
	if len(params) < cntRequiredParams {
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
		return invalidPeriod
	}

	code := m.getCurrencyCode(userID)
	rate, err := m.exchangeRateRepo.GetRate(code)
	if err != nil {
		log.Println("cannot get rate from expenseRepo:", err)
		return canNotGetRate
	}

	var sb strings.Builder
	currencyShort := getCurrencyShortByCode(code)
	report, err := m.expenseRepo.Report(m.ctx, userID, period)
	if err != nil {
		log.Println("cannot get report:", err)
		return canNotCreateReport
	}

	for _, item := range report {
		amount := float64(item.AmountInKopecks) * rate
		sb.WriteString(item.Category)
		sb.WriteString(fmt.Sprintf(": %.2f %v\n", amount/cntKopInRub, currencyShort))
	}
	return sb.String()
}

func (m *Model) getCurrencyCode(userID int64) string {
	code := m.userRepo.GetCurrency(userID)
	if code != nil {
		return *code
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
	amount, err := m.parseAmount(userID, params[1])
	if err != nil {
		log.Println("cannot parse amount:", err)
		return invalidAmount
	}
	if err := m.userRepo.SetLimit(userID, amount); err != nil {
		log.Println("cannot parse amount:", err)
		return canNotSaveLimit
	}
	return limitSaved
}

func (m *Model) delLimitHandler(userID int64) string {
	if err := m.userRepo.DelLimit(userID); err != nil {
		log.Println("cannot parse amount:", err)
		return canNotSaveLimit
	}
	return limitDeleted
}
