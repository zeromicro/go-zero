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
	t.Setenv("GOZERO_K8S_LOCAL_FALLBACK", "true")

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

func TestKubeBuilder_Build_LocalFallback(t *testing.T) {
	tests := []struct {
		name        string
		fallbackEnv string
		kubeconfig  string
		errContains string
	}{
		{
			name:        "disabled when env not set",
			fallbackEnv: "",
			errContains: "GOZERO_K8S_LOCAL_FALLBACK",
		},
		{
			name:        "disabled when env is false",
			fallbackEnv: "false",
			errContains: "GOZERO_K8S_LOCAL_FALLBACK",
		},
		{
			name:        "enabled when env is true",
			fallbackEnv: "true",
			errContains: "k8s config load failed",
		},
		{
			name:        "use custom KUBECONFIG path",
			fallbackEnv: "true",
			kubeconfig:  "/nonexistent/custom/kubeconfig",
			errContains: "/nonexistent/custom/kubeconfig",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			t.Setenv("KUBERNETES_SERVICE_HOST", "")
			t.Setenv("KUBERNETES_SERVICE_PORT", "")
			if tt.fallbackEnv != "" {
				t.Setenv("GOZERO_K8S_LOCAL_FALLBACK", tt.fallbackEnv)
			}
			if tt.kubeconfig != "" {
				t.Setenv("KUBECONFIG", tt.kubeconfig)
			}

			var b kubeBuilder
			cc := &mockedClientConn{}

			u, err := url.Parse(fmt.Sprintf("%s://my-service.default:8080", KubernetesScheme))
			assert.NoError(t, err)

			_, err = b.Build(resolver.Target{
				URL: *u,
			}, cc, resolver.BuildOptions{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}
