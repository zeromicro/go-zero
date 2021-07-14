package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/redis/redistest"
	"google.golang.org/grpc/metadata"
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
			hasError: false,
		},
		{
			name:     "strict=true",
			strict:   true,
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

	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

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
