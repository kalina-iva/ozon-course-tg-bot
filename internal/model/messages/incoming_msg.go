package messages

import (
	"fmt"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"strconv"
	"strings"
)

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Repository interface {
	NewCategory(userID int64, name string) *entity.Category
	GetCategories(userId int64) []*entity.Category
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

func (s *Model) IncomingMessage(msg Message) error {
	var params = strings.Split(msg.Text, " ")
	switch params[0] {
	case "/start":
		return s.tgClient.SendMessage("hello", msg.UserID)
	case "/newcat":
		return s.NewCatHandler(msg.UserID, params)
	case "/allcat":
		return s.AllCatHandler(msg.UserID)
	default:
		return s.tgClient.SendMessage("Неизвестная команда", msg.UserID)
	}
}

func (s *Model) NewCatHandler(userID int64, params []string) error {
	if len(params) == 1 {
		return s.tgClient.SendMessage("Нет названия категории", userID)
	}
	catName := strings.Join(params[1:], " ")
	newCat := s.repo.NewCategory(userID, catName)
	return s.tgClient.SendMessage(
		fmt.Sprintf("Создана категория %s. Ее номер %v", newCat.Name, newCat.Number),
		userID,
	)
}

func (s *Model) AllCatHandler(userID int64) error {
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

	return s.tgClient.SendMessage(text, userID)
}
