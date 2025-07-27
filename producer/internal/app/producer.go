package app

import (
	"context"
	event "github.com/ewik2k21/kafka-event-logger/common/event"
	"github.com/ewik2k21/kafka-event-logger/common/metrics"
	"github.com/ewik2k21/kafka-event-logger/producer/internal/delivery"
	"github.com/google/uuid"
	"log/slog"
	"time"
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

			eventMsg := event.NewEvent(userId.String(), eventId.String(), productId.String(), event.ToEventType(eventTypes[i%3]))
			if err := p.kafkaProducer.SendEvent(eventMsg); err != nil {
				p.logger.Error("Failed send event", slog.String("error", err.Error()))
				metrics.ProduceErrors.Inc()
				continue
			}
			metrics.EventsProduced.Inc()
			p.logger.Info("Event was send ")
			time.Sleep(4 * time.Second)
		}
	}
}
