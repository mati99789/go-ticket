package kafka

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/mati/go-ticket/internal/domain"
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	logger        *slog.Logger
	topics        []string
	msgHandler    domain.EventHandler
}

func NewKafkaConsumer(
	consumer sarama.ConsumerGroup,
	topics []string,
	logger *slog.Logger,
	event domain.EventHandler) *KafkaConsumer {
	return &KafkaConsumer{
		consumerGroup: consumer,
		topics:        topics,
		logger:        logger,
		msgHandler:    event,
	}
}

func (k *KafkaConsumer) Consume(ctx context.Context) error {
	handler := &consumerHandler{
		logger: k.logger,
		event:  k.msgHandler,
	}

	err := k.consumerGroup.Consume(ctx, k.topics, handler)
	if err != nil {
		return err
	}

	return nil
}

type consumerHandler struct {
	logger *slog.Logger
	event  domain.EventHandler
}

func (h *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Info("Setup kafka consumer")
	return nil
}

func (h *consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Info("Cleanup kafka consumer")
	return nil
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			err := h.event.Handle(session.Context(), message.Value)
			session.MarkMessage(message, "")

			if err != nil {
				h.logger.Error("Error on kafka consumer", "error", err)
				continue
			}
		}
	}
}
