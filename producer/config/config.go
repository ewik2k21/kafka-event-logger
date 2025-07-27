package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	MetricsPort  string
}

func InitConfig() *Config {
	metricsPortFlag := flag.String("metrics-port", "", "port for metrics Prometheus")
	flag.Parse()

	metricsPort := *metricsPortFlag
	if metricsPort == "" {
		metricsPort = os.Getenv("METRICS_PORT")
	}
	if metricsPort == "" {
		metricsPort = "8081"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}

	return &Config{
		KafkaBrokers: strings.Split(kafkaBrokers, ","),
		MetricsPort:  metricsPort,
	}
}
