package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/api"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/report"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

type Consumer struct {
	ctx          context.Context
	generator    *report.Generator
	reportClient api.ReportClient
}

func NewConsumerGroup(
	ctx context.Context,
	brokerList []string,
	groupID string,
	topic string,
	reportClient api.ReportClient,
	generator *report.Generator,
) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupID, config)
	if err != nil {
		return errors.Wrap(err, "cannot start consumer group")
	}

	err = consumerGroup.Consume(ctx, []string{topic}, &Consumer{
		ctx:          ctx,
		generator:    generator,
		reportClient: reportClient,
	})
	if err != nil {
		return errors.Wrap(err, "consuming via handler")
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("c - setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("c - cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var request report.Request
		err := json.Unmarshal(message.Value, &request)
		if err != nil {
			logger.Error("cannot unmarshal message", zap.Error(err))
			continue
		}

		r := c.generator.GenerateReport(c.ctx, request)
		_, err = c.reportClient.SendReport(c.ctx, &api.ReportRequest{
			UserID: request.UserID,
			Report: r,
		})
		if err != nil {
			logger.Error("cannot send generated report", zap.Error(err))
		}

		session.MarkMessage(message, "")
	}

	return nil
}
