package delivery

import (
	"github.com/IBM/sarama"
	"github.com/ewik2k21/kafka-event-logger/common/event"
	"github.com/ewik2k21/kafka-event-logger/producer/config"
	"log/slog"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	logger   *slog.Logger
}

func NewKafkaProducer(cfg *config.Config, logger *slog.Logger) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: producer, logger: logger}, nil
}

func (p *KafkaProducer) SendEvent(event event.Event) error {
	data, err := event.ToJson()
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "user-events",
		Value: sarama.ByteEncoder(data),
	}
	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
