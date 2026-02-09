// Package postgres implements the EventRepository interface using PostgreSQL.
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mati/go-ticket/internal/domain"
)

// EventRepository implements the EventRepository interface using PostgreSQL.
type EventRepository struct {
	queries *Queries
}

// NewEventRepository creates a new EventRepository.
func NewEventRepository(queries *Queries) *EventRepository {
	return &EventRepository{queries: queries}
}

// CreateEvent creates a new event in the database.
func (r *EventRepository) CreateEvent(ctx context.Context, event *domain.Event) error {
	startAt, endAt := event.StartAndEndAt()

	params := CreateEventParams{
		ID:             pgtype.UUID{Bytes: event.ID(), Valid: true},
		Name:           event.Name(),
		Price:          event.Price(),
		StartAt:        pgtype.Timestamptz{Time: startAt, Valid: true},
		EndAt:          pgtype.Timestamptz{Time: endAt, Valid: true},
		CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Capacity:       int32(event.Capacity()),
		AvailableSpots: int32(event.AvailableSpots()),
	}

	_, err := r.queries.CreateEvent(ctx, params)
	return err
}

// UpdateEvent updates an event in the database.
func (r *EventRepository) UpdateEvent(ctx context.Context, event *domain.Event) error {
	startAt, endAt := event.StartAndEndAt()

	params := UpdateEventParams{
		ID:        pgtype.UUID{Bytes: event.ID(), Valid: true},
		Name:      event.Name(),
		Price:     event.Price(),
		StartAt:   pgtype.Timestamptz{Time: startAt, Valid: true},
		EndAt:     pgtype.Timestamptz{Time: endAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Capacity:  int32(event.Capacity()),
	}

	_, err := r.queries.UpdateEvent(ctx, params)
	return err
}

// DeleteEvent deletes an event from the database.
func (r *EventRepository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteEvent(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return err
}

// GetEvent retrieves an event from the database.
func (r *EventRepository) GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	row, err := r.queries.GetEvent(ctx, pgtype.UUID{Bytes: id, Valid: true})

	if err != nil {
		return nil, domain.ErrEventNotFound
	}

	return domain.NewEventFromPersistence(
		uuid.UUID(row.ID.Bytes),
		row.Name,
		row.Price,
		row.StartAt.Time,
		row.EndAt.Time,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
		int(row.Capacity),
		int(row.AvailableSpots),
	), nil
}

// ListEvents retrieves a list of events from the database.
func (r *EventRepository) ListEvents(ctx context.Context) ([]*domain.Event, error) {
	rows, err := r.queries.ListEvents(ctx, ListEventsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, domain.ErrEventNotFound
	}

	var events []*domain.Event
	for _, row := range rows {
		event := domain.NewEventFromPersistence(
			uuid.UUID(row.ID.Bytes),
			row.Name,
			row.Price,
			row.StartAt.Time,
			row.EndAt.Time,
			row.CreatedAt.Time,
			row.UpdatedAt.Time,
			int(row.Capacity),
			int(row.AvailableSpots),
		)
		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepository) ReserveSpots(ctx context.Context, eventID uuid.UUID, spots int) error {
	_, err := r.queries.ReserveSpots(ctx, ReserveSpotsParams{
		ID:             pgtype.UUID{Bytes: eventID, Valid: true},
		AvailableSpots: int32(spots),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, errGet := r.queries.GetEvent(ctx, pgtype.UUID{Bytes: eventID, Valid: true})
			if errGet != nil {
				return domain.ErrEventNotFound
			}
			return domain.ErrEventIsFull
		}
		return err
	}
	return nil
}

func (r *EventRepository) WithTx(tx pgx.Tx) *EventRepository {
	return &EventRepository{queries: r.queries.WithTx(tx)}
}
