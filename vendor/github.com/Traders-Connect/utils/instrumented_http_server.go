package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	utils_http "github.com/Traders-Connect/utils/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// HttpSrvConfig holds instrumented http server's config
type HttpSrvConfig struct {
	ServiceName string
	APIAddr     string
	MetricsAddr string
	Logger      *zap.SugaredLogger
}

// InstrumentedHTTPServer is a generic http server includes a zap logger to be used and metrics server
type InstrumentedHTTPServer struct {
	addr          string
	Echo          *echo.Echo
	metricsServer *http.Server
	Log           *zap.SugaredLogger
}

// NewInstrumentedHTTPServer returns a new instance of the server
func NewInstrumentedHTTPServer(c HttpSrvConfig) (*InstrumentedHTTPServer, error) {
	if c.Logger == nil {
		return nil, errors.New("invalid logger")
	}

	//ignoredEndpoints := []string{"/grpc.health.v1.Health/Check"}

	// metrics
	metricsMiddleware, err := utils_http.NewMetricsMiddleware("traders_connect", c.ServiceName)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	ms := &http.Server{
		Addr:           c.MetricsAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	loggingMiddleware, err := utils_http.NewLoggingMiddleware(c.Logger)
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(loggingMiddleware.Logger)
	e.Use(metricsMiddleware.Metrics)
	e.Use(middleware.Recover())

	return &InstrumentedHTTPServer{
		addr:          c.APIAddr,
		metricsServer: ms,
		Echo:          e,
		Log:           c.Logger,
	}, nil
}

// Run starts the server
func (s *InstrumentedHTTPServer) Run(ctx context.Context, wg *sync.WaitGroup) error {
	runCtx, cancel := context.WithCancel(ctx)

	// ctx handler func
	wg.Add(1)
	go func() {
		<-runCtx.Done()
		// metrics server
		err := s.metricsServer.Shutdown(runCtx)
		if err != nil && !errors.Is(err, context.Canceled) {
			s.Log.Infow("error shutting down metrics server", "error", err.Error())
		}

		err = s.Echo.Shutdown(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			s.Log.Infow("error shutting down echo server", "error", err.Error())
		}
		s.Log.Info("shutting down echo server done gracefully")
		wg.Done()
	}()

	// starting metrics server
	s.metricsServer.RegisterOnShutdown(func() {
		s.Log.Info("metrics server shut down")
		wg.Done()
	})
	wg.Add(1)
	go func() {
		err := s.metricsServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error(err)
		}
	}()
	s.Log.Infof("starting metrics server on %s", s.metricsServer.Addr)

	s.Log.Infof("starting http server on %s", s.addr)
	wg.Add(1)
	// start http server
	go func() {
		err := s.Echo.Start(fmt.Sprintf("%s", s.addr))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Log.Error(err)
			cancel()
		}
	}()

	return nil
}
