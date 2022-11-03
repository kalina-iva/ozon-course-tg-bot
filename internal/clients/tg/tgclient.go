package tg

import (
	"context"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

type tokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
	wg     sync.WaitGroup
}

func New(tokenGetter tokenGetter) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "create NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) SendMessage(text string, cases []string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)

	if len(cases) > 0 {
		keyboard := tgbotapi.InlineKeyboardMarkup{}
		for _, value := range cases {
			var row []tgbotapi.InlineKeyboardButton
			btn := tgbotapi.NewInlineKeyboardButtonData(value, value)
			row = append(row, btn)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		}
		msg.ReplyMarkup = keyboard
	}

	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "cannot client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(ctx context.Context, msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("listening for messages")

	c.wg.Add(1)
	for update := range updates {
		if update.Message != nil {
			logger.Info(
				"starting process message",
				zap.String("from", update.Message.From.UserName),
				zap.String("msg", update.Message.Text),
			)

			startTime := time.Now()

			command, err := msgModel.IncomingMessage(ctx, messages.Message{
				Text:   update.Message.Text,
				UserID: update.Message.From.ID,
			})

			histogramResponseTime.WithLabelValues(command).Observe(time.Since(startTime).Seconds())
			if err != nil {
				logger.Error("cannot handle message", zap.Error(err))
			}
		} else if update.CallbackQuery != nil {
			logger.Info(
				"starting process callback",
				zap.String("from", update.CallbackQuery.From.UserName),
				zap.String("data", update.CallbackQuery.Data),
			)

			err := msgModel.SetCurrency(ctx, messages.CallbackQuery{
				Data:   update.CallbackQuery.Data,
				UserID: update.CallbackQuery.From.ID,
			})
			if err != nil {
				logger.Error("error processing callback", zap.Error(err))
			}
		}
	}
	c.wg.Done()
}

func (c *Client) Close() {
	c.client.StopReceivingUpdates()
	c.wg.Wait()
}
