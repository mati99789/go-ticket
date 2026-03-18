package domain

import "context"

type MessageConsumer interface {
	Consume(ctx context.Context) error
}
