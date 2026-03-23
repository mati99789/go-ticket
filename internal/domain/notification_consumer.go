package domain

import "context"

type Message struct {
	Body   []byte
	Ack    func() error
	Nack   func() error
	Reject func() error
}

type NotificationConsumer interface {
	Consume(ctx context.Context) (<-chan Message, error)
}
