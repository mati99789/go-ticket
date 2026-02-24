package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRateLimiter(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	limiter := NewRateLimiter(client, 3, 1*time.Minute)

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		allowed, err := limiter.Allow(ctx, "test")
		if err != nil {
			t.Errorf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Errorf("request %d: expected allowed, got blocked", i)
		}
	}
	allowed, err := limiter.Allow(ctx, "test")
	if err != nil {
		t.Fatalf("4th request: unexpected error: %v", err)
	}
	if allowed {
		t.Errorf("4th request: expected to be blocked, but was allowed")
	}

	mr.FastForward(time.Minute + time.Second)

	allowed, err = limiter.Allow(ctx, "test")
	if err != nil {
		t.Errorf("After window reset: unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("After window reset: expected allowed, got blocked")
	}
}
