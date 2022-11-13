package producer

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

const backoffMillisecond = 250

func NewSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
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

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot start Sarama producer")
	}

	return producer, nil
}
