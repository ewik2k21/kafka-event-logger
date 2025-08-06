package delivery

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/ewik2k21/kafka-event-logger/consumer/config"
	"log/slog"
	"time"
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	logger        *slog.Logger
}

func NewKafkaConsumer(cfg *config.Config, logger *slog.Logger) (*KafkaConsumer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaCfg.Net.DialTimeout = 10 * time.Second
	saramaCfg.Net.ReadTimeout = 10 * time.Second
	saramaCfg.Net.WriteTimeout = 10 * time.Second
	saramaCfg.Consumer.Retry.Backoff = 100

	consumerGroup, err := sarama.NewConsumerGroup(cfg.KafkaBrokers, "consumer-group", saramaCfg)
	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{consumerGroup: consumerGroup, logger: logger}, nil
}

func (c *KafkaConsumer) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	for {
		err := c.consumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			c.logger.Error("Error consume msg", slog.String("error", err.Error()))
			if ctx.Err() != nil {
				return err
			}
			time.Sleep(1 * time.Second)
			continue
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.consumerGroup.Close()
}
