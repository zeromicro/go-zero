package clientinterceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/prometheus"
	"google.golang.org/grpc"
)

func TestPromMetricInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		enable bool
		err    error
	}{
		{
			name:   "nil",
			enable: true,
			err:    nil,
		},
		{
			name:   "with error",
			enable: true,
			err:    errors.New("mock"),
		},
		{
			name: "disabled",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.enable {
				prometheus.StartAgent(prometheus.Config{
					Host: "localhost",
					Path: "/",
				})
			}
			cc := new(grpc.ClientConn)
			err := PrometheusInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}
}
