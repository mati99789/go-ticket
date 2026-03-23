package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mati/go-ticket/internal/domain"
	"github.com/redis/go-redis/v9"
)

type EmailWorker struct {
	logger      *slog.Logger
	redisClient *redis.Client
	consumer    domain.NotificationConsumer
}

func NewEmailWorker(
	logger *slog.Logger,
	redisClient *redis.Client,
	consumer domain.NotificationConsumer,
) *EmailWorker {
	return &EmailWorker{
		redisClient: redisClient,
		logger:      logger,
		consumer:    consumer,
	}
}

func (e *EmailWorker) Start(ctx context.Context) error {
	msg, err := e.consumer.Consume(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			e.logger.Info("EmailWorker shutting down")
			return nil
		case msgs := <-msg:
			var booking domain.BookingEventPayload
			err := json.Unmarshal(msgs.Body, &booking)
			if err != nil {
				e.logger.Error("failed unmarshal booking event: ", "error", err)
				_ = msgs.Reject()
				continue
			}
			key := fmt.Sprintf("email:sent:%s", booking.ID)
			exists, err := e.redisClient.Exists(ctx, key).Result()
			if err != nil {
				e.logger.Error("failed checking if email exists: ", "error", err)
				_ = msgs.Reject()
				continue
			}

			if exists > 0 {
				e.logger.Info("EmailWorker skipped due to booking exists")
				err := msgs.Ack()
				if err != nil {
					e.logger.Error("failed acknowledging email: ", "error", err)
					continue
				}
				continue
			}

			e.logger.Info("Sending email to:", "booking_id", booking.UserEmail)
			err = e.redisClient.Set(ctx, key, 1, 7*24*time.Hour).Err()
			if err != nil {
				e.logger.Error("failed sending email: ", "error", err)
				_ = msgs.Reject()
				continue
			}
			err = msgs.Ack()
			if err != nil {
				e.logger.Error("failed acknowledging email: ", "error", err)
				continue
			}
		}
	}
}
