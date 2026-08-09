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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
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

func TestKubeBuilder_BuildReturnsKubeClientError(t *testing.T) {
	restoreConfig := mockKubeConfig(t)
	defer restoreConfig()

	sentinel := errors.New("kube client failed")
	oldNewClient := newKubeClient
	defer func() {
		newKubeClient = oldNewClient
	}()
	newKubeClient = func(*rest.Config) (kubernetes.Interface, error) {
		return nil, sentinel
	}

	var b kubeBuilder
	u, err := url.Parse("k8s://dev/demo-rpc:8080")
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{URL: *u}, nil, resolver.BuildOptions{})
	assert.ErrorIs(t, err, sentinel)
}

func TestKubeBuilder_BuildReturnsAddEventHandlerError(t *testing.T) {
	restoreConfig := mockKubeConfig(t)
	defer restoreConfig()

	sentinel := errors.New("add event handler failed")
	restoreHandler := mockEndpointSliceEventHandler(t, sentinel)
	defer restoreHandler()

	var b kubeBuilder
	u, err := url.Parse("k8s://dev/demo-rpc:8080")
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{URL: *u}, nil, resolver.BuildOptions{})
	assert.ErrorIs(t, err, sentinel)
}

func TestKubeBuilder_BuildWrapsFirstEndpointSliceListError(t *testing.T) {
	restore := mockKubeClient(t)
	defer restore()

	var b kubeBuilder
	u, err := url.Parse("k8s://dev/demo-rpc")
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{URL: *u}, nil, resolver.BuildOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list EndpointSlices")
	assert.Contains(t, err.Error(), "list/watch permissions")
}

func TestKubeBuilder_BuildWrapsSecondEndpointSliceListError(t *testing.T) {
	restore := mockKubeClient(t)
	defer restore()

	var b kubeBuilder
	u, err := url.Parse("k8s://dev/demo-rpc:8080")
	assert.NoError(t, err)

	_, err = b.Build(resolver.Target{URL: *u}, nil, resolver.BuildOptions{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list EndpointSlices")
	assert.Contains(t, err.Error(), "list/watch permissions")
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
	assert.Contains(t, err.Error(), "list/watch")
}

func TestWrapEndpointSliceListErrorOther(t *testing.T) {
	svc := kube.Service{
		Namespace: "dev",
		Name:      "demo-rpc",
	}
	err := errors.New("not forbidden")
	assert.Same(t, err, wrapEndpointSliceListError(svc, err))
}

func TestEndpointSliceTweakListOptions(t *testing.T) {
	var options v1.ListOptions
	endpointSliceTweakListOptions("demo-rpc")(&options)
	assert.Equal(t, serviceSelector+"demo-rpc", options.LabelSelector)
}

func mockKubeClient(t *testing.T) func() {
	t.Helper()

	restoreConfig := mockKubeConfig(t)
	oldNewClient := newKubeClient

	newKubeClient = func(*rest.Config) (kubernetes.Interface, error) {
		cli := k8sfake.NewSimpleClientset()
		cli.PrependReactor("list", "endpointslices", func(action ktesting.Action) (bool, runtime.Object, error) {
			_ = action
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Group: "discovery.k8s.io", Resource: "endpointslices"},
				"demo-rpc",
				fmt.Errorf("forbidden"),
			)
		})
		return cli, nil
	}

	return func() {
		restoreConfig()
		newKubeClient = oldNewClient
	}
}

func mockKubeConfig(t *testing.T) func() {
	t.Helper()

	oldConfig := inClusterConfig
	inClusterConfig = func() (*rest.Config, error) {
		return &rest.Config{Host: "https://127.0.0.1"}, nil
	}

	return func() {
		inClusterConfig = oldConfig
	}
}

func mockEndpointSliceEventHandler(t *testing.T, err error) func() {
	t.Helper()

	oldAddHandler := addEndpointSliceEventHandler
	addEndpointSliceEventHandler = func(cache.SharedIndexInformer, cache.ResourceEventHandler) error {
		return err
	}

	return func() {
		addEndpointSliceEventHandler = oldAddHandler
	}
}
