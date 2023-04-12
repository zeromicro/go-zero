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
	"github.com/zeromicro/go-zero/core/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	// reset slowThreshold
	t.Cleanup(func() {
		slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
	})
	assert.Equal(t, time.Second, slowThreshold.Load())
}

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

	DontLogContentForMethod("foo")
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

func Test_shouldLogContent(t *testing.T) {
	type args struct {
		method                         string
		staticNotLoggingContentMethods []string
	}

	tests := []struct {
		name  string
		args  args
		want  bool
		setup func()
	}{
		{
			"empty",
			args{
				method: "foo",
			},
			true,
			nil,
		},
		{
			"static",
			args{
				method:                         "foo",
				staticNotLoggingContentMethods: []string{"foo"},
			},
			false,
			nil,
		},
		{
			"dynamic",
			args{
				method: "foo",
			},
			false,
			func() {
				DontLogContentForMethod("foo")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			// reset ignoreContentMethods
			t.Cleanup(func() {
				ignoreContentMethods = sync.Map{}
			})
			set := collection.NewSet()
			set.AddStr(tt.args.staticNotLoggingContentMethods...)
			assert.Equalf(t, tt.want, shouldLogContent(tt.args.method, set), "shouldLogContent(%v, %v)", tt.args.method, tt.args.staticNotLoggingContentMethods)
		})
	}
}

func Test_isSlow(t *testing.T) {
	type args struct {
		duration            time.Duration
		staticSlowThreshold time.Duration
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		setup func()
	}{
		{
			"default",
			args{
				duration: time.Millisecond * 501,
			},
			true,
			nil,
		},
		{
			"static",
			args{
				duration:            time.Millisecond * 200,
				staticSlowThreshold: time.Millisecond * 100,
			},
			true,
			nil,
		},
		{
			"dynamic",
			args{
				duration: time.Millisecond * 200,
			},
			true,
			func() {
				SetSlowThreshold(time.Millisecond * 100)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			// reset slowThreshold
			t.Cleanup(func() {
				slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
			})
			assert.Equalf(t, tt.want, isSlow(tt.args.duration, tt.args.staticSlowThreshold), "isSlow(%v, %v)", tt.args.duration, tt.args.staticSlowThreshold)
		})
	}
}
