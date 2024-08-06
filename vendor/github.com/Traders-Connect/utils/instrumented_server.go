package utils

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	utils_grpc "github.com/Traders-Connect/utils/grpc"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Config is the configuration for the API server
type Config struct {
	ServiceName      string
	APIAddr          string
	APIAddrInt       string
	MetricsAddr      string
	AuthDomain       string
	TokenAuthEnabled bool
	RpcEnpointRules  utils_grpc.RPCRules
	Logger           *zap.SugaredLogger
}

// InstrumentedServer is a generic server which includes a private gRPC endpoint, a public gRPC endpoint and a metrics endpoint.
// It also includes a zap logger to be used
type InstrumentedServer struct {
	addr               string
	addrInt            string
	GRPCServer         *grpc.Server
	GRPCServerInternal *grpc.Server
	metricsServer      *http.Server
	Log                *zap.SugaredLogger
}

// NewInstrumentedServer returns a new instance of the server
func NewInstrumentedServer(c Config) (*InstrumentedServer, error) {
	if c.Logger == nil {
		return nil, errors.New("invalid logger")
	}

	ignoredEndpoints := []string{"/grpc.health.v1.Health/Check"}

	// metrics
	metricsInterceptor, err := utils_grpc.NewGRPCMetricsInterceptor("traders_connect", c.ServiceName, ignoredEndpoints)
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

	authInterceptor, err := utils_grpc.NewGRPCAuthInterceptor(c.AuthDomain, c.RpcEnpointRules, c.TokenAuthEnabled, c.Logger)
	if err != nil {
		return nil, err
	}

	loggingInterceptor, err := utils_grpc.NewGRPCLoggingInterceptor(c.Logger, ignoredEndpoints)
	if err != nil {
		return nil, err
	}

	// grpc options
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			metricsInterceptor.Interceptor,
			authInterceptor.Interceptor,
			loggingInterceptor.Interceptor,
		),
	}

	// internal grpc options
	optsInt := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			metricsInterceptor.Interceptor,
			loggingInterceptor.Interceptor,
		),
	}

	return &InstrumentedServer{
		addr:               c.APIAddr,
		addrInt:            c.APIAddrInt,
		GRPCServer:         grpc.NewServer(opts...),
		GRPCServerInternal: grpc.NewServer(optsInt...),
		metricsServer:      ms,
		Log:                c.Logger,
	}, nil
}

// Run starts the server
func (s *InstrumentedServer) Run(ctx context.Context, wg *sync.WaitGroup) error {
	runCtx, cancel := context.WithCancel(ctx)

	// ctx handler func
	wg.Add(1)
	go func() {
		<-runCtx.Done()
		// grpc server
		s.Log.Info("stopping gRPC server")
		s.GRPCServer.Stop()
		wg.Done()
		// grpc internal server
		s.Log.Info("stopping gRPC internal server")
		s.GRPCServerInternal.Stop()
		wg.Done()

		// metrics server
		err := s.metricsServer.Shutdown(runCtx)
		if err != nil && !errors.Is(err, context.Canceled) {
			s.Log.Infow("error shutting down metrics server", "error", err.Error())
		}
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

	// starting gRPC server
	s.Log.Infof("starting gRPC server on %s", s.addr)
	sock, err := net.Listen("tcp", s.addr)
	if err != nil {
		cancel()
		return err
	}

	// starting gRPC internal server
	s.Log.Infof("starting internal gRPC server on %s", s.addrInt)
	sockInt, err := net.Listen("tcp", s.addrInt)
	if err != nil {
		cancel()
		return err
	}

	wg.Add(1)
	go func() {
		if err := s.GRPCServer.Serve(sock); err != nil {
			s.Log.Error(err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		if err := s.GRPCServerInternal.Serve(sockInt); err != nil {
			s.Log.Error(err)
			cancel()
		}
	}()

	return nil
}
