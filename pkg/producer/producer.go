package producer

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

const backoffMillisecond = 250

type Producer struct {
	syncProducer sarama.SyncProducer
	reportTopic  string
}

func NewSyncProducer(brokerList []string, reportTopic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = backoffMillisecond * time.Millisecond
	if config.Producer.Idempotent {
		config.Producer.Retry.Max = 1
		config.Net.MaxOpenRequests = 1
	}
	config.Producer.Return.Successes = true
	_ = config.Producer.Partitioner

	p, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot start Sarama syncProducer")
	}

	return &Producer{
		syncProducer: p,
		reportTopic:  reportTopic,
	}, nil
}

func (p *Producer) SendReportMessage(userID int64, msg []byte) error {
	_, _, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: p.reportTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("report%d", userID)),
		Value: sarama.StringEncoder(msg),
	})
	if err != nil {
		return errors.Wrap(err, "cannot produce report message")
	}
	return nil
}
