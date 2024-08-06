package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Check performs a healthcheck
func (s *NotificationService) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	//err := s.Db.PingContext(ctx)
	//if err != nil {
	//	s.Log.Errorw("healthcheck error", "dependency", "database", "error", err)
	//	return &grpc_health_v1.HealthCheckResponse{
	//		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	//	}, err
	//}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch performs a stream healthcheck
func (s *NotificationService) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "unimplemented")
}
