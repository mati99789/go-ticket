package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRateLimiter(client *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	cmd := r.client.Incr(ctx, key)
	count, err := cmd.Result()

	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := r.client.Expire(ctx, key, r.window).Err(); err != nil {
			return false, err
		}
	}

	return count <= int64(r.limit), nil
}
