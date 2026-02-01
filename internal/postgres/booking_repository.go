package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mati/go-ticket/internal/domain"
)

type BookingRepository struct {
	Queries *Queries
}

func NewBookingRepository(queries *Queries) *BookingRepository {
	return &BookingRepository{
		Queries: queries,
	}
}

func (br *BookingRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	params := CreateBookingParams{
		ID:        pgtype.UUID{Bytes: booking.ID(), Valid: true},
		EventID:   pgtype.UUID{Bytes: booking.EventID(), Valid: true},
		UserEmail: booking.UserEmail(),
		Status:    string(booking.Status()),
		CreatedAt: pgtype.Timestamptz{Time: booking.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: booking.UpdatedAt(), Valid: true},
	}
	_, err := br.Queries.CreateBooking(ctx, params)
	return err
}

func (br *BookingRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	row, err := br.Queries.GetBookingByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return domain.UnmarshalBooking(
		uuid.UUID(row.ID.Bytes),
		uuid.UUID(row.EventID.Bytes),
		row.UserEmail,
		domain.BookingStatus(row.Status),
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}

func (br *BookingRepository) UpdateBooking(ctx context.Context, booking *domain.Booking) error {
	params := UpdateBookingParams{
		ID:        pgtype.UUID{Bytes: booking.ID(), Valid: true},
		EventID:   pgtype.UUID{Bytes: booking.EventID(), Valid: true},
		UserEmail: booking.UserEmail(),
		Status:    string(booking.Status()),
		UpdatedAt: pgtype.Timestamptz{Time: booking.UpdatedAt(), Valid: true},
	}
	_, err := br.Queries.UpdateBooking(ctx, params)
	return err
}

func (br *BookingRepository) DeleteBooking(ctx context.Context, id uuid.UUID) error {
	err := br.Queries.DeleteBooking(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return err
}

func (br *BookingRepository) ConfirmBooking(ctx context.Context, id uuid.UUID) error {
	_, err := br.Queries.ConfirmBooking(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return err
}

func (br *BookingRepository) CancelBooking(ctx context.Context, id uuid.UUID) error {
	_, err := br.Queries.CancelBooking(ctx, pgtype.UUID{Bytes: id, Valid: true})
	return err
}

func (br *BookingRepository) ListBookings(ctx context.Context) ([]domain.Booking, error) {
	rows, err := br.Queries.ListBookings(ctx, ListBookingsParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}
	var bookings []domain.Booking
	for _, row := range rows {
		booking := domain.UnmarshalBooking(
			uuid.UUID(row.ID.Bytes),
			uuid.UUID(row.EventID.Bytes),
			row.UserEmail,
			domain.BookingStatus(row.Status),
			row.CreatedAt.Time,
			row.UpdatedAt.Time,
		)
		bookings = append(bookings, *booking)
	}
	return bookings, nil
}

func (r *BookingRepository) WithTx(tx pgx.Tx) *BookingRepository {
	return &BookingRepository{Queries: r.Queries.WithTx(tx)}
}
