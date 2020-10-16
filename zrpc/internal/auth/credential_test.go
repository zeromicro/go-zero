package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestParseCredential(t *testing.T) {
	tests := []struct {
		name        string
		withNil     bool
		withEmptyMd bool
		app         string
		token       string
	}{
		{
			name:    "nil",
			withNil: true,
		},
		{
			name:        "empty md",
			withEmptyMd: true,
		},
		{
			name: "empty",
		},
		{
			name:  "valid",
			app:   "foo",
			token: "bar",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			if test.withNil {
				ctx = context.Background()
			} else if test.withEmptyMd {
				ctx = metadata.NewIncomingContext(context.Background(), metadata.MD{})
			} else {
				md := metadata.New(map[string]string{
					"app":   test.app,
					"token": test.token,
				})
				ctx = metadata.NewIncomingContext(context.Background(), md)
			}
			cred := ParseCredential(ctx)
			assert.False(t, cred.RequireTransportSecurity())
			m, err := cred.GetRequestMetadata(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, test.app, m[appKey])
			assert.Equal(t, test.token, m[tokenKey])
		})
	}
}
