package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mati/go-ticket/internal/domain"
)

type OutBoxRepository struct {
	queries *Queries
}

func NewOutBoxRepository(queries *Queries) *OutBoxRepository {
	return &OutBoxRepository{queries: queries}
}

func (r *OutBoxRepository) Create(ctx context.Context, event *domain.OutboxEvent) error {
	params := CreateOutboxEventParams{
		ID:          pgtype.UUID{Bytes: event.ID(), Valid: true},
		EventName:   event.EventName(),
		EventData:   event.EventData(),
		Destination: event.Destination(),
	}
	_, err := r.queries.CreateOutboxEvent(ctx, params)
	return err
}

func (r *OutBoxRepository) GetPendingEvents(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	dbEvents, err := r.queries.GetPendingOutboxEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}
	var events []*domain.OutboxEvent
	for _, dbEvent := range dbEvents {
		event := domain.ReconstructOutboxEvent(
			uuid.UUID(dbEvent.ID.Bytes),
			dbEvent.EventName,
			dbEvent.EventData,
			dbEvent.Status.String,
			dbEvent.CreatedAt.Time,
			dbEvent.UpdatedAt.Time,
			dbEvent.Destination,
		)
		events = append(events, event)
	}
	return events, nil

}

func (r *OutBoxRepository) MarkAsProcessed(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = r.queries.MarkOutBoxEventAsProcessed(ctx, pgtype.UUID{Bytes: parsedID, Valid: true})
	if err != nil {
		return err
	}

	return nil
}

func (r *OutBoxRepository) WithTx(tx pgx.Tx) *OutBoxRepository {
	return &OutBoxRepository{
		queries: r.queries.WithTx(tx),
	}
}
