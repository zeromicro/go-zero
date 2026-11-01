package serverinterceptors

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/zrpc/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestStreamAuthorizeInterceptor(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		token    string
		hasError bool
	}{
		{
			name:     "unknown app",
			hasError: true,
		},
		{
			name:     "valid token",
			app:      "foo",
			token:    "bar",
			hasError: false,
		},
		{
			name:     "wrong token",
			app:      "foo",
			token:    "error",
			hasError: true,
		},
	}

	store := redistest.CreateRedis(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.app) > 0 {
				assert.Nil(t, store.Hset("apps", test.app, test.token))
				defer store.Hdel("apps", test.app)
			}

			authenticator, err := auth.NewAuthenticator(store, "apps")
			assert.Nil(t, err)
			interceptor := StreamAuthorizeInterceptor(authenticator)
			md := metadata.New(map[string]string{
				"app":   "foo",
				"token": "bar",
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			stream := mockedStream{ctx: ctx}
			called := false
			err = interceptor(nil, stream, nil, func(_ any, _ grpc.ServerStream) error {
				called = true
				return nil
			})
			assert.Equal(t, !test.hasError, called)
			if test.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUnaryAuthorizeInterceptor(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		token    string
		hasError bool
	}{
		{
			name:     "unknown app",
			hasError: true,
		},
		{
			name:     "valid token",
			app:      "foo",
			token:    "bar",
			hasError: false,
		},
		{
			name:     "wrong token",
			app:      "foo",
			token:    "error",
			hasError: true,
		},
	}

	store := redistest.CreateRedis(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.app) > 0 {
				assert.Nil(t, store.Hset("apps", test.app, test.token))
				defer store.Hdel("apps", test.app)
			}

			authenticator, err := auth.NewAuthenticator(store, "apps")
			assert.Nil(t, err)
			interceptor := UnaryAuthorizeInterceptor(authenticator)
			md := metadata.New(map[string]string{
				"app":   "foo",
				"token": "bar",
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			called := false
			_, err = interceptor(ctx, nil, nil,
				func(ctx context.Context, req any) (any, error) {
					called = true
					return nil, nil
				})
			assert.Equal(t, !test.hasError, called)
			if test.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			_, err = interceptor(context.Background(), nil, nil,
				func(ctx context.Context, req any) (any, error) {
					return nil, nil
				})
			assert.NotNil(t, err)

			md = nil
			ctx = metadata.NewIncomingContext(context.Background(), md)
			_, err = interceptor(ctx, nil, nil,
				func(ctx context.Context, req any) (any, error) {
					return nil, nil
				})
			assert.NotNil(t, err)

			md = metadata.New(map[string]string{
				"app":   "",
				"token": "",
			})
			ctx = metadata.NewIncomingContext(context.Background(), md)
			_, err = interceptor(ctx, nil, nil,
				func(ctx context.Context, req any) (any, error) {
					return nil, nil
				})
			assert.NotNil(t, err)
		})
	}
}

type mockedStream struct {
	ctx context.Context
}

func (m mockedStream) SetHeader(_ metadata.MD) error {
	return nil
}

func (m mockedStream) SendHeader(_ metadata.MD) error {
	return nil
}

func (m mockedStream) SetTrailer(_ metadata.MD) {
}

func (m mockedStream) Context() context.Context {
	return m.ctx
}

func (m mockedStream) SendMsg(_ any) error {
	return nil
}

func (m mockedStream) RecvMsg(_ any) error {
	return nil
}

func TestAuthorizeInterceptorsRejectLookupErrors(t *testing.T) {
	for _, backendError := range []bool{false, true} {
		t.Run(fmt.Sprintf("backendError=%v", backendError), func(t *testing.T) {
			mr := miniredis.RunT(t)
			if backendError {
				mr.SetError("ERR injected backend failure")
			}
			authenticator, err := auth.NewAuthenticator(redis.New(mr.Addr()), "apps")
			require.NoError(t, err)
			ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("app", "unknown", "token", "arbitrary"))
			called := false
			_, err = UnaryAuthorizeInterceptor(authenticator)(ctx, nil, nil, func(context.Context, any) (any, error) {
				called = true
				return nil, nil
			})
			assert.Equal(t, codes.Unauthenticated, status.Code(err))
			assert.False(t, called)
			err = StreamAuthorizeInterceptor(authenticator)(nil, mockedStream{ctx: ctx}, nil, func(any, grpc.ServerStream) error {
				called = true
				return nil
			})
			assert.Equal(t, codes.Unauthenticated, status.Code(err))
			assert.False(t, called)
		})
	}
}
