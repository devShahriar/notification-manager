package grpc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GRPCLoggingInterceptor struct {
	logger           *zap.SugaredLogger
	ignoredEndpoints map[string]struct{}
}

func NewGRPCLoggingInterceptor(logger *zap.SugaredLogger, ignoredEndpoints []string) (*GRPCLoggingInterceptor, error) {
	if logger == nil {
		return nil, fmt.Errorf("invalid logger provided")
	}

	ie := make(map[string]struct{})
	for _, endpoint := range ignoredEndpoints {
		ie[endpoint] = struct{}{}
	}

	return &GRPCLoggingInterceptor{logger: logger, ignoredEndpoints: ie}, nil
}

func (i *GRPCLoggingInterceptor) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := i.ignoredEndpoints[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	startTime := time.Now()
	resp, err := handler(ctx, req)
	statusCode := status.Code(err)
	duration := time.Since(startTime)
	userId := GetUserID(ctx)

	i.logger.Infow("received request", "rpc", info.FullMethod, "status_code", statusCode.String(), "user", userId, "duration", duration)

	return resp, err
}
