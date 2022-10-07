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
	cntKopInRub     = 100
	defaultCurrency = "RUB"
)

type currencyRepository interface {
	GetRate(currency string) (float64, error)
}

type messageSender interface {
	SendMessage(text string, cases []string, userID int64) (int, error)
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

func (s *Model) IncomingMessage(msg Message) (err error) {
	var text string
	var cases []string

	params := strings.Split(msg.Text, " ")
	switch params[0] {
	case "/start":
		text = manual
	case "/help":
		text = manual
	case "/newexpense":
		text = s.newExpenseHandler(msg.UserID, params)
	case "/report":
		text = s.reportHandler(msg.UserID, params)
	case "/setcurrency":
		text = chooseCurrency
		cases = []string{"RUB", "USD", "EUR", "CNY"}
	default:
		text = unknownCommand
	}
	_, err = s.tgClient.SendMessage(text, cases, msg.UserID)

	return
}

func (s *Model) DoAction(msg CallbackQuery) (err error) {
	s.repo.SetCurrency(msg.UserID, msg.Data)
	_, err = s.tgClient.SendMessage(currencySaved, nil, msg.UserID)
	return
}

func (s *Model) newExpenseHandler(userID int64, params []string) string {
	const cntRequiredParams = 3
	if len(params) < cntRequiredParams {
		return needCategoryAndAmount
	}
	category := params[1]
	amount, err := s.parseAmount(userID, params[2])
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

	s.repo.NewExpense(userID, category, amount, date)

	return expenseAdded
}

func (s *Model) parseAmount(userID int64, amountStr string) (uint64, error) {
	const bitSize = 64
	amount, err := strconv.ParseFloat(amountStr, bitSize)
	if err != nil {
		return 0, errors.Wrap(err, "parse amount")
	}
	if amount <= 0 {
		return 0, errInvalidAmount
	}
	rate, err := s.getCustomCurrencyRate(userID)
	if err != nil {
		return 0, errors.Wrap(err, "get currency rate")
	}
	amount /= rate
	return uint64(math.Round(amount * cntKopInRub)), nil
}

func (s *Model) reportHandler(userID int64, params []string) string {
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

	rate, err := s.getCustomCurrencyRate(userID)
	if err != nil {
		return canNotGateRate
	}

	report := s.repo.NewReport(userID, period)
	var sb strings.Builder
	for _, item := range report {
		amount := float64(item.AmountInKopecks) * rate
		sb.WriteString(item.Category)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%.2f\n", amount/cntKopInRub))
	}
	return sb.String()
}

func (s *Model) getCustomCurrencyRate(userID int64) (float64, error) {
	currency := s.repo.GetCurrency(userID)
	if len(currency) == 0 {
		currency = defaultCurrency
	}

	rate, err := s.currencyRepository.GetRate(currency)
	if err != nil {
		return 0, err
	}
	return rate, nil
}
