package messages

import (
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/mocks/messages"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	sender.EXPECT().SendMessage("hello", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	sender.EXPECT().SendMessage("Неизвестная команда", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewCatCommand_ShouldAnswerWithError(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	sender.EXPECT().SendMessage("Нет названия категории", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newcat",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewCatCommand_NameOfSeveralWords(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	repo.EXPECT().NewCategory(int64(123), "new category").Return(&entity.Category{
		Number: 1,
		Name:   "new category",
	})
	sender.EXPECT().SendMessage("Создана категория new category. Ее номер 1", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newcat new category",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAllCatCommand_NoCategories(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	repo.EXPECT().GetCategories(int64(123)).Return(nil)
	sender.EXPECT().SendMessage("Создана категория new category. Ее номер 1", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/allcat",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnAllCatCommand_TwoCategories(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	var categories []*entity.Category
	categories = append(
		categories,
		&entity.Category{Number: 1, Name: "new category"},
		&entity.Category{Number: 2, Name: "second category"},
	)
	repo.EXPECT().GetCategories(int64(123)).Return(categories)
	sender.EXPECT().SendMessage("1. new category\n2. second category\n", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/allcat",
		UserID: 123,
	})

	assert.NoError(t, err)
}
