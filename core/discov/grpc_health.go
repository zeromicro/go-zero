package discov

import (
	"context"

	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Server for gRPC health checking protocol.
type gRPCHealthServer struct {
	metrics *stat.Metrics
}

// Registers grpc health check protocol to *grpc.Server
func RegisterGRPCHealthCheck(server *grpc.Server, metrics *stat.Metrics) {
	grpc_health_v1.RegisterHealthServer(server, newGRPCHealthServer(metrics))
}

// Returns new server that is responsible for GRPC Health Checking Protocol.
func newGRPCHealthServer(metrics *stat.Metrics) grpc_health_v1.HealthServer {
	return &gRPCHealthServer{
		metrics: metrics,
	}
}

// Endpoint for grpc-health-probe or kubernetes 1.23 probes to query for unary status.
func (s *gRPCHealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	checkService := req.GetService()
	if checkService != "" {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_UNKNOWN,
		}, status.Error(codes.Unimplemented, "individual service statuses is unavailable")
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Endpoint for kubernetes 1.23 probes to query for streaming status. This is NOT implemented for grpc-health-probe.
func (*gRPCHealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, watch grpc_health_v1.Health_WatchServer) error {
	return watch.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
