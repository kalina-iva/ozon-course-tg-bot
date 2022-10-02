package messages

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"strconv"
	"strings"
	"time"
)

type messageSender interface {
	SendMessage(text string, userID int64) error
}

type repository interface {
	NewCategory(userID int64, name string) *entity.Category
	GetCategories(userID int64) []*entity.Category
	NewExpense(userID int64, category entity.Category, amount float64, date int64)
	NewReport(userID int64, period int64) []*entity.Report
}

type Model struct {
	tgClient messageSender
	repo     repository
}

func New(tgClient messageSender, repo repository) *Model {
	return &Model{
		tgClient: tgClient,
		repo:     repo,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) IncomingMessage(msg Message) error {
	var text string

	var params = strings.Split(msg.Text, " ")
	switch params[0] {
	case "/start":
		text = manual
	case "/help":
		text = manual
	case "/newcat":
		text = s.newCatHandler(msg.UserID, params)
	case "/allcat":
		text = s.allCatHandler(msg.UserID)
	case "/newexpense":
		text = s.newExpenseHandler(msg.UserID, params)
	case "/report":
		text = s.reportHandler(msg.UserID, params)
	default:
		text = unknownCommand
	}
	return s.tgClient.SendMessage(text, msg.UserID)
}

func (s *Model) newCatHandler(userID int64, params []string) string {
	if len(params) == 1 {
		return "Нет названия категории"
	}
	catName := strings.Join(params[1:], " ")
	newCat := s.repo.NewCategory(userID, catName)
	return fmt.Sprintf("Создана категория %s. Ее номер %v", newCat.Name, newCat.Number)
}

func (s *Model) allCatHandler(userID int64) string {
	categories := s.repo.GetCategories(userID)
	var text string
	if len(categories) == 0 {
		text = noCategories
	} else {
		var sb strings.Builder
		for _, cat := range categories {
			sb.WriteString(strconv.Itoa(cat.Number))
			sb.WriteString(". ")
			sb.WriteString(cat.Name)
			sb.WriteString("\n")
		}
		text = sb.String()
	}

	return text
}

func (s *Model) newExpenseHandler(userID int64, params []string) string {
	if len(params) < 3 {
		return needCategoryAndAmount
	}
	category, err := s.checkCategory(userID, params[1])
	if err != nil {
		if errors.Is(err, invalidCategoryNumberErr) {
			return invalidCategoryNumber
		}
		if errors.Is(err, categoryNotFoundErr) {
			return "Не найдена категория с номером " + params[1]
		}
	}
	amount, err := s.checkAmount(params[2])
	if err != nil {
		return invalidAmount
	}
	var date int64
	if len(params) == 4 {
		t, err := time.Parse("01-02-2006", params[3])
		if err != nil {
			return invalidDate
		}
		date = t.Unix()
	} else {
		date = time.Now().Unix()
	}

	s.repo.NewExpense(userID, *category, amount, date)

	return expenseAdded
}

func (s *Model) checkCategory(userID int64, categoryNumber string) (*entity.Category, error) {
	number, err := strconv.Atoi(categoryNumber)
	if err != nil {
		return nil, invalidCategoryNumberErr
	}
	categories := s.repo.GetCategories(userID)
	for _, cat := range categories {
		if cat.Number == number {
			return cat, nil
		}
	}
	return nil, categoryNotFoundErr
}

func (s *Model) checkAmount(amountStr string) (float64, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}
	if amount <= 0 {
		return 0, errors.New("amount cannot be negative or 0")
	}
	return amount, nil
}

func (s *Model) reportHandler(userID int64, params []string) string {
	if len(params) < 2 {
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

	report := s.repo.NewReport(userID, period)
	var sb strings.Builder
	for _, item := range report {
		sb.WriteString(item.Category.Name)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%.2f\n", item.Amount))
	}
	return sb.String()
}
