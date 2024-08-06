package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsMiddleware struct {
	count    *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func NewMetricsMiddleware(namespace, subsystem string) (*MetricsMiddleware, error) {
	if namespace == "" || subsystem == "" {
		return nil, fmt.Errorf("empty namespace and/or subsystem")
	}

	var (
		total = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_requests_total",
				Help:      "The total number of http requests per function",
			},
			[]string{"http", "status"},
		)
		duration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_requests_duration_seconds",
				Help:      "The amount of time http request functions take",
			},
			[]string{"http", "status"},
		)
	)

	return &MetricsMiddleware{
		count:    total,
		duration: duration,
	}, nil
}

func (m *MetricsMiddleware) Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		duration := time.Since(startTime)

		m.count.WithLabelValues(c.Path(), strconv.FormatInt(int64(c.Response().Status), 10)).Inc()
		m.duration.WithLabelValues(c.Path(), strconv.FormatInt(int64(c.Response().Status), 10)).Observe(duration.Seconds())
		return nil
	}
}
