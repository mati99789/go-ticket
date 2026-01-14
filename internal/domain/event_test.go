package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/domain"
)

//nolint:funlen
func TestNewEvent(t *testing.T) {
	type args struct {
		id      uuid.UUID
		name    string
		price   int64
		startAt time.Time
		endAt   time.Time
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid event",
			args: args{
				id:      uuid.New(),
				name:    "Metallica Concert",
				price:   10000, // 100.00
				startAt: time.Now().Add(24 * time.Hour),
				endAt:   time.Now().Add(26 * time.Hour),
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			args: args{
				id:      uuid.New(),
				name:    "",
				price:   100,
				startAt: time.Now(),
				endAt:   time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrEventNameEmpty,
		},
		{
			name: "negative price",
			args: args{
				id:      uuid.New(),
				name:    "Metallica Concert",
				price:   -100,
				startAt: time.Now(),
				endAt:   time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrEventPriceNegative,
		},
		{
			name: "endAt before startAt",
			args: args{
				id:      uuid.New(),
				name:    "Metallica Concert",
				price:   100,
				startAt: time.Now().Add(time.Hour),
				endAt:   time.Now(),
			},
			wantErr: domain.ErrEventStartAfterEnd,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewEvent(tt.args.id, tt.args.name, tt.args.price, tt.args.startAt, tt.args.endAt)

			if err != tt.wantErr {
				t.Errorf("NewEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil {
				return
			}

			if got == nil {
				t.Errorf("NewEvent() got = nil")
				return
			}

			if got.Name() != tt.args.name {
				t.Errorf("NewEvent() name = %v, want %v", got.Name(), tt.args.name)
			}
		})
	}
}
