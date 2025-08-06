package server

import (
	"context"
	"github.com/ewik2k21/kafka-event-logger/common/metrics"
	"github.com/ewik2k21/kafka-event-logger/consumer/config"
	"github.com/ewik2k21/kafka-event-logger/consumer/internal/app"
	"github.com/ewik2k21/kafka-event-logger/consumer/internal/delivery"
	"os"
	"os/signal"
	"syscall"

	"log/slog"
	"sync"
)

func Execute(ctx context.Context, cfg *config.Config, logger *slog.Logger) {
	wg := sync.WaitGroup{}

	metrics.InitMetrics(ctx, &wg, cfg.PushgatewayURL, "consumer", logger)

	kafkaConsumer, err := delivery.NewKafkaConsumer(cfg, logger)
	if err != nil {
		logger.Error("Failed to init kafka consumer", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer kafkaConsumer.Close()

	consumerApp, err := app.NewConsumer(cfg, kafkaConsumer, logger)
	if err != nil {
		logger.Error("Failed init consumer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumerApp.Run(ctx)
	}()

	logger.Info("Consumer start")

	<-sigChan
	logger.Info("Signal to shutdown")
	cancel()
	wg.Wait()
	logger.Info("Consumer stopped")
}
