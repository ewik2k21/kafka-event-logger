package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers   []string
	PushgatewayURL string
	EventsLogPath  string
}

func NewConfig() Config {

	pushgatewayURLFlag := flag.String("pushgateway-url", "", "URL для Pushgateway Prometheus")
	flag.Parse()

	pushgatewayURL := *pushgatewayURLFlag
	if pushgatewayURL == "" {
		pushgatewayURL = os.Getenv("PUSHGATEWAY_URL")
	}
	if pushgatewayURL == "" {
		pushgatewayURL = "http://pushgateway:9091"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}

	eventsLogPath := os.Getenv("EVENTS_LOG_PATH")
	if eventsLogPath == "" {
		eventsLogPath = "/data/events.log"
	}

	return Config{
		KafkaBrokers:   strings.Split(kafkaBrokers, ","),
		PushgatewayURL: pushgatewayURL,
		EventsLogPath:  eventsLogPath,
	}
}
