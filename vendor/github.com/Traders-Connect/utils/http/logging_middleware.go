package http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *zap.SugaredLogger
}

func NewLoggingMiddleware(logger *zap.SugaredLogger) (*LoggingMiddleware, error) {
	if logger == nil {
		return nil, fmt.Errorf("invalid logger provided")
	}
	return &LoggingMiddleware{logger: logger}, nil
}

func (m *LoggingMiddleware) Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		duration := time.Since(startTime)

		m.logger.Infow("received request", "endpoint", c.Path(), "status_code", strconv.FormatInt(int64(c.Response().Status), 10), "duration", duration)
		return nil
	}
}
