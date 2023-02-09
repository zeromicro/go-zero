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
	var b kubeBuilder
	u, err := url.Parse(fmt.Sprintf("%s://%s", KubernetesScheme, "a,b"))
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{
		URL: *u,
	}, nil, resolver.BuildOptions{})
	assert.Error(t, err)

	u, err = url.Parse(fmt.Sprintf("%s://%s:9100/a:b:c", KubernetesScheme, "a,b,c,d"))
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{
		URL: *u,
	}, nil, resolver.BuildOptions{})
	assert.Error(t, err)
}

func TestIsNonBlockNotFoundErr(t *testing.T) {

	tests := []struct {
		name     string
		nonBlock bool
		err      error
		hasErr   bool
	}{
		{
			name:     "block true , endpoints not found error ",
			nonBlock: true,
			err:      fmt.Errorf("endpoints app-rpc-svc not found"),
			hasErr:   false,
		},
		{
			name:     "block true , other error ",
			nonBlock: true,
			err:      fmt.Errorf("other error"),
			hasErr:   true,
		},
		{
			name:     "block false , endpoints app-rpc-svc not found ",
			nonBlock: false,
			err:      fmt.Errorf("other error"),
			hasErr:   true,
		},
		{
			name:     "block false , endpoints app-rpc-svc not found ",
			nonBlock: false,
			err:      fmt.Errorf("other error"),
			hasErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkEndpointsErr(test.nonBlock, test.err)
			if test.hasErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
