package domain

import "context"

type MessageBroker interface {
	Publish(ctx context.Context, event *OutboxEvent) error
}
