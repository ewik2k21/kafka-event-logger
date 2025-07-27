package app

import (
	"context"
	event "github.com/ewik2k21/kafka-event-logger/common/event"
	"github.com/ewik2k21/kafka-event-logger/producer/internal/delivery"
	"github.com/google/uuid"
	"log/slog"
)

type Producer struct {
	kafkaProducer *delivery.KafkaProducer
	logger        *slog.Logger
}

func NewProducer(producer *delivery.KafkaProducer, logger *slog.Logger) *Producer {
	return &Producer{kafkaProducer: producer, logger: logger}
}

func (p *Producer) Run(ctx context.Context) {
	eventTypes := []string{"view", "add_to_cart", "purchase"}
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			p.logger.Info("Producer finished work")
			return
		default:
			userId := uuid.New()
			eventId := uuid.New()
			productId := uuid.New()
			eventType := i % 3
			event := event.NewEvent(userId.String(), eventId.String(), productId.String(), event.)

		}
	}
}
