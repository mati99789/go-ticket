package workers

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/mati/go-ticket/internal/domain"
)

type KafkaConsumerWorker struct {
	consumer domain.MessageConsumer
	logger   *slog.Logger
}

func NewKafkaConsumerWorker(consumer domain.MessageConsumer, logger *slog.Logger) *KafkaConsumerWorker {
	return &KafkaConsumerWorker{
		consumer: consumer,
		logger:   logger,
	}
}

func (k *KafkaConsumerWorker) Start(ctx context.Context) error {
	k.logger.Info("Starting Kafka Consumer Worker")
	const (
		initialBackoff = time.Second
		maxBackoff     = 10 * time.Second
	)
	backoff := initialBackoff
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		err := k.consumer.Consume(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			k.logger.Error("Error consuming Kafka Message", "error", err)

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(backoff):
			}

			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		} else {
			backoff = initialBackoff
		}
	}
}
