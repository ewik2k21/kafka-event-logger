package app

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ewik2k21/kafka-event-logger/common/event"
	"github.com/ewik2k21/kafka-event-logger/common/metrics"
	"github.com/ewik2k21/kafka-event-logger/consumer/config"
	"github.com/ewik2k21/kafka-event-logger/consumer/internal/delivery"
	"log/slog"
	"os"
)

type Consumer struct {
	cfg           config.Config
	kafkaConsumer *delivery.KafkaConsumer
	logger        *slog.Logger
}

func NewConsumer(cfg *config.Config, kafkaConsumer *delivery.KafkaConsumer, logger *slog.Logger) (*Consumer, error) {
	return &Consumer{
		cfg:           *cfg,
		kafkaConsumer: kafkaConsumer,
		logger:        logger,
	}, nil
}

type consumerHandler struct {
	logger        *slog.Logger
	eventsLogPath string
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Run(ctx context.Context) {
	handler := &consumerHandler{
		logger:        c.logger,
		eventsLogPath: c.cfg.EventsLogPath,
	}
	topics := []string{"user-events"}

	go func() {
		for {
			if err := c.kafkaConsumer.Consume(ctx, topics, handler); err != nil {
				c.logger.Error("Failed consume in for ", slog.String("error", err.Error()))
				if ctx.Err() != nil {
					return
				}
			}
			c.logger.Info("Message consumed")
		}
	}()
}

func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var evt event.Event
		if err := json.Unmarshal(message.Value, &evt); err != nil {
			h.logger.Error("error from unmarshal event", slog.String("error", err.Error()))
			metrics.ConsumedErrors.Inc()
			continue
		}

		file, err := os.OpenFile(h.eventsLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			h.logger.Error("error open file for logs", slog.String("error", err.Error()))
			metrics.ConsumedErrors.Inc()
			return err
		}
		defer file.Close()

		eventJSON, err := json.Marshal(evt)
		if err != nil {
			h.logger.Error("error from marshal event")
			metrics.ConsumedErrors.Inc()
			return err
		}

		if _, err = file.WriteString(string(eventJSON) + "\n"); err != nil {
			h.logger.Error("failed write to file", slog.String("error", err.Error()))
			metrics.ConsumedErrors.Inc()
			return err
		}

		metrics.EventsConsumed.Inc()
		session.MarkMessage(message, "")
	}
	return nil
}
