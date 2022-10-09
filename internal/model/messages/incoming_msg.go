package messages

import (
	"fmt"
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

type currencyRepository interface {
	GetRate(code string) (float64, error)
}

type messageSender interface {
	SendMessage(text string, cases []string, userID int64) error
}

type repository interface {
	NewExpense(userID int64, category string, amount uint64, date int64)
	NewReport(userID int64, period int64) []*entity.Report
	SetCurrency(userID int64, currency string)
	GetCurrency(userID int64) string
}

type Model struct {
	tgClient           messageSender
	repo               repository
	currencyRepository currencyRepository
}

func New(tgClient messageSender, repo repository, currencyRepo currencyRepository) *Model {
	return &Model{
		tgClient:           tgClient,
		repo:               repo,
		currencyRepository: currencyRepo,
	}
}

type Message struct {
	Text   string
	UserID int64
}

type CallbackQuery struct {
	ID     int
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
	default:
		text = unknownCommand
	}
	return m.tgClient.SendMessage(text, cases, msg.UserID)
}

func (m *Model) SetCurrency(msg CallbackQuery) (err error) {
	m.repo.SetCurrency(msg.UserID, msg.Data)
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
		return invalidAmount
	}

	var date int64
	if len(params) == cntRequiredParams+1 {
		t, err := time.Parse("01-02-2006", params[3])
		if err != nil {
			return invalidDate
		}
		date = t.Unix()
	} else {
		date = time.Now().Unix()
	}

	m.repo.NewExpense(userID, category, amount, date)

	return expenseAdded
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
	rate, err := m.currencyRepository.GetRate(code)
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
	var period int64
	now := time.Now()
	switch params[1] {
	case "y":
		period = now.AddDate(-1, 0, 0).Unix()
	case "m":
		period = now.AddDate(0, -1, 0).Unix()
	case "w":
		period = now.AddDate(0, 0, -7).Unix()
	default:
		return invalidPeriod
	}

	code := m.getCurrencyCode(userID)
	rate, err := m.currencyRepository.GetRate(code)
	if err != nil {
		return canNotGetRate
	}

	currencyShort := getCurrencyShortByCode(code)
	report := m.repo.NewReport(userID, period)
	var sb strings.Builder
	for _, item := range report {
		amount := float64(item.AmountInKopecks) * rate
		sb.WriteString(item.Category)
		sb.WriteString(fmt.Sprintf(": %.2f %v\n", amount/cntKopInRub, currencyShort))
	}
	return sb.String()
}

func (m *Model) getCurrencyCode(userID int64) string {
	code := m.repo.GetCurrency(userID)
	if len(code) == 0 {
		code = DefaultCurrencyCode
	}
	return code
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
