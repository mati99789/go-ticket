package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/domain"
)

func TestNewBooking(t *testing.T) {
	type args struct {
		id        uuid.UUID
		eventID   uuid.UUID
		userEmail string
		status    domain.BookingStatus
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid booking",
			args: args{
				id:        uuid.New(),
				eventID:   uuid.New(),
				userEmail: "user@example.com",
				status:    domain.BookingStatusPending,
			},
			wantErr: nil,
		},
		{
			name: "nil booking ID",
			args: args{
				id:        uuid.Nil,
				eventID:   uuid.New(),
				userEmail: "user@example.com",
				status:    domain.BookingStatusPending,
			},
			wantErr: domain.ErrBookingIDNil,
		},
		{
			name: "nil event ID",
			args: args{
				id:        uuid.New(),
				eventID:   uuid.Nil,
				userEmail: "user@example.com",
				status:    domain.BookingStatusPending,
			},
			wantErr: domain.ErrBookingEventIDInvalid,
		},
		{
			name: "empty user email",
			args: args{
				id:        uuid.New(),
				eventID:   uuid.New(),
				userEmail: "",
				status:    domain.BookingStatusPending,
			},
			wantErr: domain.ErrBookingUserEmailEmpty,
		},
		{
			name: "invalid status",
			args: args{
				id:        uuid.New(),
				eventID:   uuid.New(),
				userEmail: "user@example.com",
				status:    "invalid_status",
			},
			wantErr: domain.ErrBookingStatusInvalid,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewBooking(tt.args.id, tt.args.eventID, tt.args.userEmail, tt.args.status)

			if err != tt.wantErr {
				t.Errorf("NewBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				return
			}

			if got == nil {
				t.Errorf("NewBooking() got = nil")
				return
			}

			if got.ID() != tt.args.id {
				t.Errorf("NewBooking() ID = %v, want %v", got.ID(), tt.args.id)
			}

			if got.EventID() != tt.args.eventID {
				t.Errorf("NewBooking() EventID = %v, want %v", got.EventID(), tt.args.eventID)
			}

			if got.UserEmail() != tt.args.userEmail {
				t.Errorf("NewBooking() UserEmail = %v, want %v", got.UserEmail(), tt.args.userEmail)
			}

			if got.Status() != tt.args.status {
				t.Errorf("NewBooking() Status = %v, want %v", got.Status(), tt.args.status)
			}

		})
	}
}
