package workers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mati/go-ticket/internal/domain"
)

type OutboxRelay struct {
	outboxRepo domain.OutboxRepository
	broker     domain.MessageBroker
}

func NewOutboxRelay(outboxRepo domain.OutboxRepository, broker domain.MessageBroker) *OutboxRelay {
	return &OutboxRelay{outboxRepo: outboxRepo, broker: broker}
}

func (r *OutboxRelay) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 2)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Outbox Relay is shuttiing down...")
			return nil
		case <-ticker.C:
			err := r.processOutbox(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (r *OutboxRelay) processOutbox(ctx context.Context) error {
	events, err := r.outboxRepo.GetPendingEvents(ctx, 50)

	if len(events) == 0 {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to fetch outbox events: %w", err)
	}

	for _, event := range events {
		err := r.broker.Publish(ctx, event)
		if err != nil {
			slog.Error("Failed to publish event", "error", err)
			continue
		}
		if err := r.outboxRepo.MarkAsProcessed(ctx, event.ID().String()); err != nil {
			slog.Error("Failed to mark event as processed", "error", err)
			continue
		}
	}

	return nil
}
