package interfaces

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type healthcheckHandler struct{}

func NewHealthServer() grpc_health_v1.HealthServer {
	return &healthcheckHandler{}
}

func (s *healthcheckHandler) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *healthcheckHandler) Watch(req *grpc_health_v1.HealthCheckRequest, ws grpc_health_v1.Health_WatchServer) error {
	return nil
}
