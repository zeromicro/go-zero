package serverinterceptors

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func TestUnaryStatInterceptor(t *testing.T) {
	metrics := stat.NewMetrics("mock")
	interceptor := UnaryStatInterceptor(metrics, StatConf{})
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
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
		req      any
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
				logDuration(test.ctx, "foo", test.req, test.duration,
					collection.NewSet(), 0)
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
		req      any
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
	// reset ignoreContentMethods
	t.Cleanup(func() {
		ignoreContentMethods = sync.Map{}
	})
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.NotPanics(t, func() {
				logDuration(test.ctx, "foo", test.req, test.duration,
					collection.NewSet(), 0)
			})
		})
	}
}
