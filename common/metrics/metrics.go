package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"log/slog"
	"sync"
	"time"
)

var (
	EventsProduced = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_produced_total",
			Help: "Общее кол-во отправленных событий",
		})
	EventsConsumed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "events_consumed_total",
			Help: "Общее кол-во принятых событий",
		})
	ProduceErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "produce_errors_total",
			Help: "Общее кол-во ошибок при отправке сообщений",
		})
	ConsumedErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "consume_errors_total",
			Help: "Общее кол-во ошибок при получение сообщений",
		})
)

func InitMetrics(ctx context.Context, wg *sync.WaitGroup, pushgatewayURL, jobName string, logger *slog.Logger) {
	prometheus.MustRegister(EventsProduced, EventsConsumed, ProduceErrors, ConsumedErrors)

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		sendMetrics := func() error {
			return push.New(pushgatewayURL, jobName).
				Collector(EventsProduced).
				Collector(EventsConsumed).
				Collector(ProduceErrors).
				Collector(ConsumedErrors).
				Push()
		}

		for {
			select {
			case <-ctx.Done():
				if err := sendMetrics(); err != nil {
					logger.Error("failed last metrics push", slog.String("error", err.Error()))
				}
				return
			case <-ticker.C:
				if err := sendMetrics(); err != nil {
					logger.Error("failed push metrics")
				}
			}
		}
	}()
}
