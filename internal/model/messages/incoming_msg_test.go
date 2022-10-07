package messages

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/mocks/messages"
)

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	sender.EXPECT().SendMessage(`Привет! Это дневник расходов.
Описание команд:
/newexpense {category} {amount} {date} - добавление нового расхода. Если дата не указана, используется текущая
/report {y|m|w} - получение отчета за последний год/месяц/неделю
`, nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithHelpMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	sender.EXPECT().SendMessage("Неизвестная команда", nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_WrongAmount(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	sender.EXPECT().SendMessage("Некорректная сумма расхода", nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense category 0 02-10-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_incorrectDate(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	sender.EXPECT().SendMessage("Некорректная дата", nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense category 76.10 29-02-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnNewExpenseCommand_onOk(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	repo.EXPECT().NewExpense(int64(123), "category", uint64(7610), int64(1644451200))
	sender.EXPECT().SendMessage("Расход добавлен", nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/newexpense category 76.10 02-10-2022",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnReportCommand_withoutPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	sender.EXPECT().SendMessage("Необходимо указать период", nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/report",
		UserID: 123,
	})

	assert.NoError(t, err)
}

func Test_OnReportCommand_wrongPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockmessageSender(ctrl)
	repo := mocks.NewMockrepository(ctrl)
	currencyRepo := mocks.NewMockcurrencyRepository(ctrl)
	model := New(sender, repo, currencyRepo)

	sender.EXPECT().SendMessage("Некорректный период", nil, int64(123))

	err := model.IncomingMessage(Message{
		Text:   "/report a",
		UserID: 123,
	})

	assert.NoError(t, err)
}
