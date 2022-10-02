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

	categories := []*entity.Category{
		{Number: 1, Name: "new category"},
		{Number: 2, Name: "second category"},
	}
	repo.EXPECT().GetCategories(int64(123)).Return(categories)
	sender.EXPECT().SendMessage("1. new category\n2. second category\n", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/allcat",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_CategoryNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	repo.EXPECT().GetCategories(int64(123)).Return(nil)
	sender.EXPECT().SendMessage("Не найдена категория с номером 1", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense 1 76.10 02-10-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_WrongCategory(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	sender.EXPECT().SendMessage("Некорректный номер категории", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense cat 76.10 02-10-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_WrongAmount(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	categories := []*entity.Category{
		{Number: 1, Name: "new category"},
	}
	repo.EXPECT().GetCategories(int64(123)).Return(categories)
	sender.EXPECT().SendMessage("Некорректная сумма расхода", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense 1 0 02-10-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_onOk(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	repo := mocks.NewMockRepository(ctrl)
	model := New(sender, repo)

	categories := []*entity.Category{
		{Number: 1, Name: "new category"},
	}
	repo.EXPECT().GetCategories(int64(123)).Return(categories)
	repo.EXPECT().NewExpense(int64(123), 1, 76.10, int64(1644451200))
	sender.EXPECT().SendMessage("Расход добавлен", int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense 1 76.10 02-10-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}
