package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mati/go-ticket/internal/domain"
)

type KafkaBroker struct {
	producer sarama.SyncProducer
}

func NewKafkaBroker(producer sarama.SyncProducer) *KafkaBroker {
	return &KafkaBroker{producer: producer}
}

func (b *KafkaBroker) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	msg := &sarama.ProducerMessage{
		Topic: event.Destination(),
		Key:   sarama.StringEncoder(event.AggregateID().String()),
		Value: sarama.ByteEncoder(event.EventData()),
	}

	_, _, err := b.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
