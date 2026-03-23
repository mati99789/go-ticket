package rabbitmq

import (
	"context"
	"fmt"

	"github.com/mati/go-ticket/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMqConsumer(conn *amqp.Connection, queueName string) (*RabbitMQConsumer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQConsumer{
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (c *RabbitMQConsumer) Consume(ctx context.Context) (<-chan domain.Message, error) {
	deliveries, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume: %w", err)
	}

	msgs := make(chan domain.Message)
	go func() {
		for delivery := range deliveries {
			msgs <- domain.Message{
				Body:   delivery.Body,
				Ack:    func() error { return delivery.Ack(false) },
				Nack:   func() error { return delivery.Nack(false, false) },
				Reject: func() error { return delivery.Reject(false) },
			}
		}
	}()

	return msgs, nil
}

func (c *RabbitMQConsumer) Close() error {
	return c.channel.Close()
}
