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
	resyncInterval = 5 * time.Minute
	nameSelector   = "metadata.name="
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
		// getting endpoints is only to get the port
		endpoints, err := cs.CoreV1().Endpoints(svc.Namespace).Get(
			context.Background(), svc.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}

		svc.Port = int(endpoints.Subsets[0].Ports[0].Port)
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
			options.FieldSelector = nameSelector + svc.Name
		}))
	in := inf.Core().V1().Endpoints()
	_, err = in.Informer().AddEventHandler(handler)
	if err != nil {
		return nil, err
	}

	// get the initial endpoints, cannot use the previous endpoints,
	// because the endpoints may be updated before/after the informer is started.
	endpoints, err := cs.CoreV1().Endpoints(svc.Namespace).Get(
		context.Background(), svc.Name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	handler.Update(endpoints)

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
