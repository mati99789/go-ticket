package workers

import (
	"context"
	"log/slog"

	"github.com/mati/go-ticket/internal/domain"
)

type LogBroker struct{}

func NewLogBroker() *LogBroker {
	return &LogBroker{}
}

func (b *LogBroker) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	slog.Info("MOCK BROKER: Publishing event to the world!",
		"id", event.ID(),
		"type", event.EventName(),
		"destination", event.Destination())
	return nil
}
