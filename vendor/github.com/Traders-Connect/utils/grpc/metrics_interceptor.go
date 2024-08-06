package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GRPCMetricsInterceptor struct {
	count            *prometheus.CounterVec
	duration         *prometheus.HistogramVec
	ignoredEndpoints map[string]struct{}
}

func NewGRPCMetricsInterceptor(namespace, subsystem string, ignoredEndpoints []string) (*GRPCMetricsInterceptor, error) {
	if namespace == "" || subsystem == "" {
		return nil, fmt.Errorf("empty namespace and/or subsystem")
	}

	var (
		grpc_total = promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "grpc_requests_total",
				Help:      "The total number of grpc requests per function",
			},
			[]string{"rpc", "status"},
		)
		grpc_duration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "grpc_requests_duration_seconds",
				Help:      "The amount of time grpc request functions take",
			},
			[]string{"rpc", "status"},
		)
	)

	ie := make(map[string]struct{})
	for _, endpoint := range ignoredEndpoints {
		ie[endpoint] = struct{}{}
	}

	return &GRPCMetricsInterceptor{
		count:            grpc_total,
		duration:         grpc_duration,
		ignoredEndpoints: ie,
	}, nil
}

func (i *GRPCMetricsInterceptor) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := i.ignoredEndpoints[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	startTime := time.Now()
	resp, err := handler(ctx, req)
	statusCode := status.Code(err)
	duration := time.Since(startTime)

	i.count.WithLabelValues(info.FullMethod, statusCode.String()).Inc()
	i.duration.WithLabelValues(info.FullMethod, statusCode.String()).Observe(duration.Seconds())

	return resp, err
}
