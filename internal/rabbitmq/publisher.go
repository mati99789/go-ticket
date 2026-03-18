package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mati/go-ticket/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

func NewRabbitMQPublisher(conn *amqp.Connection) (*RabbitMQPublisher, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	exchangeName := "booking.events"
	err = channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false, false, false, nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	queueName := "booking.notifications"
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	routingKey := "booking.created"
	err = channel.QueueBind(
		queue.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &RabbitMQPublisher{
		channel:    channel,
		exchange:   exchangeName,
		routingKey: routingKey,
	}, nil
}

func (r *RabbitMQPublisher) Publish(ctx context.Context, booking *domain.BookingNotification) error {
	body, err := json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("failed to marshal booking: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		r.exchange,
		r.routingKey,
		false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}

func (r *RabbitMQPublisher) Close() error {
	return r.channel.Close()
}
