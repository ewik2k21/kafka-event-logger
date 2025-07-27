package server

import (
	"context"
	"github.com/ewik2k21/kafka-event-logger/common/metrics"
	"github.com/ewik2k21/kafka-event-logger/producer/config"
	"github.com/ewik2k21/kafka-event-logger/producer/internal/app"
	"github.com/ewik2k21/kafka-event-logger/producer/internal/delivery"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Execute(ctx context.Context, cfg *config.Config, logger *slog.Logger) {
	ms := metrics.InitMetrics(cfg.MetricsPort, logger)

	kafkaProducer, err := delivery.NewKafkaProducer(cfg, logger)
	if err != nil {
		logger.Error("Failed to init kafka producer", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer kafkaProducer.Close()

	producerApp := app.NewProducer(kafkaProducer, logger)

	ctx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		producerApp.Run(ctx)
	}()

	<-sigChan
	logger.Info("Signal to shutdown")
	if err := ms.Shutdown(ctx, logger); err != nil {
		logger.Error("error shutdown metrics ", slog.String("error", err.Error()))
	}
	logger.Info("Metrics server shutdown")
	cancel()
	wg.Wait()
	logger.Info("Producer shutdown")

}
