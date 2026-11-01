package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthenticator(t *testing.T) {
	tests := []struct {
		name     string
		app      string
		token    string
		strict   bool
		hasError bool
	}{
		{
			name:     "strict=false",
			strict:   false,
			hasError: true,
		},
		{
			name:     "strict=true",
			strict:   true,
			hasError: true,
		},
		{
			name:     "strict=false,with token",
			app:      "foo",
			token:    "bar",
			hasError: false,
		},
		{
			name:     "strict=false,with error token",
			app:      "foo",
			token:    "error",
			hasError: true,
		},
		{
			name:     "strict=true,with token",
			app:      "foo",
			token:    "bar",
			strict:   true,
			hasError: false,
		},
		{
			name:     "strict=true,with error token",
			app:      "foo",
			token:    "error",
			strict:   true,
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

			authenticator, err := NewAuthenticator(store, "apps", test.strict)
			assert.Nil(t, err)
			assert.NotNil(t, authenticator.Authenticate(context.Background()))
			md := metadata.New(map[string]string{})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			assert.NotNil(t, authenticator.Authenticate(ctx))
			md = metadata.New(map[string]string{
				"app":   "",
				"token": "",
			})
			ctx = metadata.NewIncomingContext(context.Background(), md)
			assert.NotNil(t, authenticator.Authenticate(ctx))
			md = metadata.New(map[string]string{
				"app":   "foo",
				"token": "bar",
			})
			ctx = metadata.NewIncomingContext(context.Background(), md)
			err = authenticator.Authenticate(ctx)
			if test.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthenticatorFailClosed(t *testing.T) {
	for _, strict := range []bool{false, true} {
		for _, scenario := range []string{"unknown", "backend_error", "wrong_type", "valid", "wrong_token", "empty_token"} {
			t.Run(fmt.Sprintf("strict=%v/%s", strict, scenario), func(t *testing.T) {
				mr := miniredis.RunT(t)
				if scenario == "backend_error" {
					mr.SetError("ERR injected backend failure")
				} else if scenario == "wrong_type" {
					mr.Set("apps", "not-a-hash")
				} else if scenario != "unknown" {
					mr.HSet("apps", "foo", "bar")
				}
				a, err := NewAuthenticator(redis.New(mr.Addr()), "apps", strict)
				if err != nil {
					t.Fatal(err)
				}
				token := "bar"
				if scenario == "wrong_token" {
					token = "wrong"
				}
				if scenario == "empty_token" {
					token = ""
				}
				ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("app", "foo", "token", token))
				err = a.Authenticate(ctx)
				want := codes.Unauthenticated
				if scenario == "valid" {
					want = codes.OK
				}
				if got := status.Code(err); got != want {
					t.Fatalf("got %v, want %v", got, want)
				}
			})
		}
	}
}
