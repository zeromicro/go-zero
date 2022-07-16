package serverinterceptors

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func TestUnaryStatInterceptor(t *testing.T) {
	metrics := stat.NewMetrics("mock")
	interceptor := UnaryStatInterceptor(metrics)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestLogDuration(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	assert.Nil(t, err)
	assert.True(t, len(addrs) > 0)

	tests := []struct {
		name     string
		ctx      context.Context
		req      interface{}
		duration time.Duration
	}{
		{
			name: "normal",
			ctx:  context.Background(),
			req:  "foo",
		},
		{
			name: "bad req",
			ctx:  context.Background(),
			req:  make(chan lang.PlaceholderType), // not marshalable
		},
		{
			name:     "timeout",
			ctx:      context.Background(),
			req:      "foo",
			duration: time.Second,
		},
		{
			name: "timeout",
			ctx: peer.NewContext(context.Background(), &peer.Peer{
				Addr: addrs[0],
			}),
			req: "foo",
		},
		{
			name:     "timeout",
			ctx:      context.Background(),
			req:      "foo",
			duration: slowThreshold.Load() + time.Second,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.NotPanics(t, func() {
				logDuration(test.ctx, "foo", test.req, test.duration)
			})
		})
	}
}

func TestLogDurationWithoutContent(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	assert.Nil(t, err)
	assert.True(t, len(addrs) > 0)

	tests := []struct {
		name     string
		ctx      context.Context
		req      interface{}
		duration time.Duration
	}{
		{
			name: "normal",
			ctx:  context.Background(),
			req:  "foo",
		},
		{
			name: "bad req",
			ctx:  context.Background(),
			req:  make(chan lang.PlaceholderType), // not marshalable
		},
		{
			name:     "timeout",
			ctx:      context.Background(),
			req:      "foo",
			duration: time.Second,
		},
		{
			name: "timeout",
			ctx: peer.NewContext(context.Background(), &peer.Peer{
				Addr: addrs[0],
			}),
			req: "foo",
		},
		{
			name:     "timeout",
			ctx:      context.Background(),
			req:      "foo",
			duration: slowThreshold.Load() + time.Second,
		},
	}

	DontLogContentForMethod("foo")
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.NotPanics(t, func() {
				logDuration(test.ctx, "foo", test.req, test.duration)
			})
		})
	}
}
