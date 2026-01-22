package internal

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestKubeBuilder_Scheme(t *testing.T) {
	var b kubeBuilder
	assert.Equal(t, KubernetesScheme, b.Scheme())
}

func TestKubeBuilder_Build(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("KUBERNETES_SERVICE_HOST", "")
	t.Setenv("KUBERNETES_SERVICE_PORT", "")

	var b kubeBuilder
	cc := &mockedClientConn{}

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid host",
			input: fmt.Sprintf("%s://%s", KubernetesScheme, "a,b"),
		},
		{
			name:  "bad endpoint format",
			input: fmt.Sprintf("%s://%s:9100/a:b:c", KubernetesScheme, "a,b,c,d"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.input)
			assert.NoError(t, err)

			_, err = b.Build(resolver.Target{
				URL: *u,
			}, cc, resolver.BuildOptions{})
			assert.Error(t, err)
		})
	}
}