package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/prometheus"
	"google.golang.org/grpc"
)

func TestUnaryPromMetricInterceptor_Disabled(t *testing.T) {
	_, err := UnaryPrometheusInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryPromMetricInterceptor_Enabled(t *testing.T) {
	prometheus.StartAgent(prometheus.Config{
		Host: "localhost",
		Path: "/",
	})
	_, err := UnaryPrometheusInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestSetRpcServerReqDurBuckets(t *testing.T) {
	SetRpcServerReqDurBuckets([]float64{0.1})
	assert.Equal(t, []float64{0.1}, rpcServerReqDurBuckets)
}
