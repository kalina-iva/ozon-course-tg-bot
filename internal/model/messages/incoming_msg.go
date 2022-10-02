package messages

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"strconv"
	"strings"
	"time"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Repository interface {
	NewCategory(userId int64, name string) *entity.Category
	GetCategories(userId int64) []*entity.Category
	NewExpense(userId int64, category entity.Category, amount float64, date int64)
	NewReport(userId int64, period int64) []*entity.Report
}

type Model struct {
	tgClient MessageSender
	repo     Repository
}

func New(tgClient MessageSender, repo Repository) *Model {
	return &Model{
		tgClient: tgClient,
		repo:     repo,
	}
}

type Message struct {
	Text   string
	UserID int64
}

var CategoryNotFound = errors.New("category not found")
var InvalidCategoryNumber = errors.New("invalid category number")
var Manual = `Привет! Это дневник расходов.
Описание команд:
/newcat {name} - добавление новой категории
/allcat - просмотр всех категории
/newexpense {categoryNumber} {amount} {date} - добавление нового расхода. Если дата не указана, используется текущая
/report {y|m|w} - получение отчета за последний год/месяц/неделю
`

func (s *Model) IncomingMessage(msg Message) error {
	var text string

	var params = strings.Split(msg.Text, " ")
	switch params[0] {
	case "/start":
		text = Manual
	case "/help":
		text = Manual
	case "/newcat":
		text = s.newCatHandler(msg.UserID, params)
	case "/allcat":
		text = s.allCatHandler(msg.UserID)
	case "/newexpense":
		text = s.newExpenseHandler(msg.UserID, params)
	case "/report":
		text = s.reportHandler(msg.UserID, params)
	default:
		text = "Неизвестная команда"
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
		text = "Нет категорий"
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

func (s *Model) newExpenseHandler(userId int64, params []string) string {
	if len(params) < 3 {
		return "Необходимо указать категорию и сумму"
	}
	category, err := s.checkCategory(userId, params[1])
	if err != nil {
		if errors.Is(err, InvalidCategoryNumber) {
			return "Некорректный номер категории"
		}
		if errors.Is(err, CategoryNotFound) {
			return "Не найдена категория с номером " + params[1]
		}
	}
	amount, err := s.checkAmount(params[2])
	if err != nil {
		return "Некорректная сумма расхода"
	}
	var date int64
	if len(params) == 4 {
		t, err := time.Parse("01-02-2006", params[3])
		if err != nil {
			return "Некорректная дата"
		}
		date = t.Unix()
	} else {
		date = time.Now().Unix()
	}

	s.repo.NewExpense(userId, *category, amount, date)

	return "Расход добавлен"
}

func (s *Model) checkCategory(userId int64, categoryNumber string) (*entity.Category, error) {
	number, err := strconv.Atoi(categoryNumber)
	if err != nil {
		return nil, InvalidCategoryNumber
	}
	categories := s.repo.GetCategories(userId)
	for _, cat := range categories {
		if cat.Number == number {
			return cat, nil
		}
	}
	return nil, CategoryNotFound
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

func (s *Model) reportHandler(userId int64, params []string) string {
	if len(params) < 2 {
		return "Необходимо указать период"
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
		return "Некорректный период"
	}

	report := s.repo.NewReport(userId, period)
	var sb strings.Builder
	for _, item := range report {
		sb.WriteString(item.Category.Name)
		sb.WriteString(": ")
		sb.WriteString(fmt.Sprintf("%.2f\n", item.Amount))
	}
	return sb.String()
}
