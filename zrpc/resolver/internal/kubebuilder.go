//go:build !no_k8s

package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc/resolver/internal/kube"
	"google.golang.org/grpc/resolver"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	resyncInterval  = 5 * time.Minute
	serviceSelector = "kubernetes.io/service-name="
)

type kubeResolver struct {
	cc     resolver.ClientConn
	inf    informers.SharedInformerFactory
	stopCh chan struct{}
}

func (r *kubeResolver) Close() {
	close(r.stopCh)
}

func (r *kubeResolver) ResolveNow(_ resolver.ResolveNowOptions) {}

func (r *kubeResolver) start() {
	threading.GoSafe(func() {
		r.inf.Start(r.stopCh)
	})
}

type kubeBuilder struct{}

func (b *kubeBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	_ resolver.BuildOptions) (resolver.Resolver, error) {
	svc, err := kube.ParseTarget(target)
	if err != nil {
		return nil, err
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if svc.Port == 0 {
		endpointSlices, err := cs.DiscoveryV1().EndpointSlices(svc.Namespace).List(context.Background(), v1.ListOptions{
			LabelSelector: serviceSelector + svc.Name,
		})
		if err != nil {
			return nil, err
		}
		if len(endpointSlices.Items) == 0 {
			return nil, fmt.Errorf("no endpoint slices found for service %s in namespace %s", svc.Name, svc.Namespace)
		}

		// Since this resolver is used for in-cluster service discovery,
		// there should be at least one port available.
		svc.Port = int(*endpointSlices.Items[0].Ports[0].Port)
	}

	handler := kube.NewEventHandler(func(endpoints []string) {
		endpoints = subset(endpoints, subsetSize)
		addrs := make([]resolver.Address, 0, len(endpoints))
		for _, val := range endpoints {
			addrs = append(addrs, resolver.Address{
				Addr: fmt.Sprintf("%s:%d", val, svc.Port),
			})
		}

		if err := cc.UpdateState(resolver.State{
			Addresses: addrs,
		}); err != nil {
			logx.Error(err)
		}
	})
	inf := informers.NewSharedInformerFactoryWithOptions(cs, resyncInterval,
		informers.WithNamespace(svc.Namespace),
		informers.WithTweakListOptions(func(options *v1.ListOptions) {
			options.LabelSelector = serviceSelector + svc.Name
		}))
	in := inf.Discovery().V1().EndpointSlices()
	_, err = in.Informer().AddEventHandler(handler)
	if err != nil {
		return nil, err
	}

	// get the initial endpoint slices, cannot use the previous endpoint slices,
	// because the endpoint slices may be updated before/after the informer is started.
	endpointSlices, err := cs.DiscoveryV1().EndpointSlices(svc.Namespace).List(
		context.Background(), v1.ListOptions{
			LabelSelector: serviceSelector + svc.Name,
		})
	if err != nil {
		return nil, err
	}

	for _, endpointSlice := range endpointSlices.Items {
		handler.Update(&endpointSlice)
	}

	r := &kubeResolver{
		cc:     cc,
		inf:    inf,
		stopCh: make(chan struct{}),
	}
	r.start()

	return r, nil
}

func (b *kubeBuilder) Scheme() string {
	return KubernetesScheme
}
