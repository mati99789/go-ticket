package postgres

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestEventRepository_CreateEvent(t *testing.T) {
	ctx := context.Background()

	pool := setupDb(ctx, t)

	//Create event
	event := createTestEvent(ctx, t, pool)

	//Get event
	retrieved := getEventFromDB(ctx, t, pool, event.ID())

	assert.Equal(t, event.ID(), retrieved.ID())
	assert.Equal(t, event.Name(), retrieved.Name())
	assert.Equal(t, event.Price(), retrieved.Price())
}

func TestEventRepository_GetEvent_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := setupDb(ctx, t)

	fakeID := uuid.New()
	queries := New(pool)
	eventRepository := NewEventRepository(queries)

	//Get event
	_, err := eventRepository.GetEvent(ctx, fakeID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEventNotFound)
}

func TestEventRepository_ReserveSpots_Success(t *testing.T) {
	ctx := context.Background()
	pool := setupDb(ctx, t)

	//Create event
	event := createTestEvent(ctx, t, pool, withCapacity(100))

	queries := New(pool)
	eventRepository := NewEventRepository(queries)

	//Reserve spots
	err := eventRepository.ReserveSpots(ctx, event.ID(), 10)

	retrieved := getEventFromDB(ctx, t, pool, event.ID())

	assert.NoError(t, err)
	assert.Equal(t, 90, retrieved.AvailableSpots())
}

func TestEventRepository_ReserveSpots_NotEnough(t *testing.T) {
	ctx := context.Background()
	pool := setupDb(ctx, t)

	//Create event
	event := createTestEvent(ctx, t, pool, withCapacity(5))

	queries := New(pool)
	eventRepository := NewEventRepository(queries)

	//Reserve spots
	err := eventRepository.ReserveSpots(ctx, event.ID(), 10)

	assert.Error(t, err)

	retrieved := getEventFromDB(ctx, t, pool, event.ID())
	assert.Equal(t, retrieved.AvailableSpots(), 5)
	assert.ErrorIs(t, err, domain.ErrEventIsFull)
}

func TestEventRepository_ReserveSpots_Concurrent(t *testing.T) {
	ctx := context.Background()
	pool := setupDb(ctx, t)

	//Create event
	event := createTestEvent(ctx, t, pool, withCapacity(100))

	queries := New(pool)
	eventRepository := NewEventRepository(queries)

	numGorutines := 200
	var wg sync.WaitGroup
	successCount := atomic.Int32{}

	for i := 0; i < numGorutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := eventRepository.ReserveSpots(ctx, event.ID(), 1)
			if err == nil {
				successCount.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, successCount.Load(), int32(100))

	retrieved := getEventFromDB(ctx, t, pool, event.ID())

	assert.Equal(t, retrieved.AvailableSpots(), 0)
}
