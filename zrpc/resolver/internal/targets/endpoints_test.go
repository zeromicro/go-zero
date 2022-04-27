package targets

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestGetAuthority(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "test",
			url:  "direct://my_authority/localhost",
			want: "my_authority",
		},
		{
			name: "test with port",
			url:  "direct://my_authority/localhost:8080",
			want: "my_authority",
		},
		{
			name: "test with multiple hosts",
			url:  "direct://my_authority1,my_authority2/localhost,localhost",
			want: "my_authority1,my_authority2",
		},
		{
			name: "test with multiple hosts with port",
			url:  "direct://my_authority1:3000,my_authority2:3001/localhost:8080,localhost:8081",
			want: "my_authority1:3000,my_authority2:3001",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := url.Parse(test.url)
			assert.Nil(t, err)
			target := resolver.Target{
				URL: *uri,
			}
			assert.Equal(t, test.want, GetAuthority(target))
		})
	}
}

func TestGetEndpoints(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "test",
			url:  "direct:///localhost",
			want: "localhost",
		},
		{
			name: "test with port",
			url:  "direct:///localhost:8080",
			want: "localhost:8080",
		},
		{
			name: "test with multiple hosts",
			url:  "direct:///localhost,localhost",
			want: "localhost,localhost",
		},
		{
			name: "test with multiple hosts with port",
			url:  "direct:///localhost:8080,localhost:8081",
			want: "localhost:8080,localhost:8081",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := url.Parse(test.url)
			assert.Nil(t, err)
			target := resolver.Target{
				URL: *uri,
			}
			assert.Equal(t, test.want, GetEndpoints(target))
		})
	}
}
