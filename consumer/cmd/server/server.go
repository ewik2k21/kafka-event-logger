package server

import (
	"context"
	"github.com/ewik2k21/kafka-event-logger/common/metrics"
	"github.com/ewik2k21/kafka-event-logger/config/config"
	"log/slog"
	"sync"
)

func Execute(ctx context.Context, cfg *config.Config, logger *slog.Logger) {
	wg := sync.WaitGroup{}

	metrics.InitMetrics(ctx, &wg, cfg.PushGatewayURL, "consumer", logger)

}
