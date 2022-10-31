package tg

import (
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages"
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

func (c *Client) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	log.Println("listening for messages")

	c.wg.Add(1)
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			err := msgModel.IncomingMessage(messages.Message{
				Text:   update.Message.Text,
				UserID: update.Message.From.ID,
			})
			if err != nil {
				log.Println("error processing message:", err)
			}
		} else if update.CallbackQuery != nil {
			log.Printf("[%s] callback %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)

			err := msgModel.SetCurrency(messages.CallbackQuery{
				Data:   update.CallbackQuery.Data,
				UserID: update.CallbackQuery.From.ID,
			})
			if err != nil {
				log.Println("error processing callback:", err)
			}
		}
	}
	c.wg.Done()
}

func (c *Client) Close() {
	c.client.StopReceivingUpdates()
	c.wg.Wait()
}