package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers   []string
	PushGatewayURL string
}

func InitConfig() *Config {
	pushGatewayURLFlag := flag.String("pushgateway-url", "", "url for metrics Prometheus")
	flag.Parse()

	pushGatewayURL := *pushGatewayURLFlag
	if pushGatewayURL == "" {
		pushGatewayURL = os.Getenv("PUSHGATEWAY_URL")
	}
	if pushGatewayURL == "" {
		pushGatewayURL = "http://pushgateway:9091"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka1:9092"
	}
	fmt.Println(kafkaBrokers)
	return &Config{
		KafkaBrokers:   strings.Split(kafkaBrokers, ","),
		PushGatewayURL: pushGatewayURL,
	}
}
