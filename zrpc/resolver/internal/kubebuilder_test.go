package internal

import (
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/zrpc/resolver/internal/kube"
	"google.golang.org/grpc/resolver"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func TestWrapEndpointSliceListError(t *testing.T) {
	svc := kube.Service{
		Namespace: "dev",
		Name:      "demo-rpc",
	}
	forbidden := apierrors.NewForbidden(
		schema.GroupResource{Group: "discovery.k8s.io", Resource: "endpointslices"},
		svc.Name,
		fmt.Errorf("forbidden"),
	)

	err := wrapEndpointSliceListError(svc, forbidden)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "EndpointSlices")
	assert.Contains(t, err.Error(), "endpointslices.discovery.k8s.io")
	assert.Contains(t, err.Error(), "get/list/watch")
}

func TestWrapEndpointSliceListErrorOther(t *testing.T) {
	svc := kube.Service{
		Namespace: "dev",
		Name:      "demo-rpc",
	}
	err := errors.New("not forbidden")
	assert.Same(t, err, wrapEndpointSliceListError(svc, err))
}
