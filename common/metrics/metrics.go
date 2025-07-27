package metrics

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
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

type MetricsServer struct {
	server *http.Server
}

func InitMetrics(port string, logger *slog.Logger) *MetricsServer {
	prometheus.MustRegister(EventsProduced, EventsConsumed, ProduceErrors, ConsumedErrors)

	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ms := &MetricsServer{server: server}

	go func() {
		logger.Info("Start metrics server on ", slog.String("port", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start metrics server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	return ms
}

func (ms *MetricsServer) Shutdown(ctx context.Context, logger *slog.Logger) error {
	logger.Info("Server shutting down ... ")
	return ms.server.Shutdown(ctx)
}
