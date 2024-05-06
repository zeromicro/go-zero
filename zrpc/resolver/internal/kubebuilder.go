package internal

import (
	"context"
	"fmt"
	"runtime/debug"
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
	stopCh chan struct{}
	inf    informers.SharedInformerFactory
}

func (r *kubeResolver) start() {
	threading.GoSafe(func() {
		r.inf.Start(r.stopCh)
	})
}

func (r *kubeResolver) ResolveNow(_ resolver.ResolveNowOptions) {}

func (r *kubeResolver) Close() {
	close(r.stopCh)
}

type kubeBuilder struct{}

func (b *kubeBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	_ resolver.BuildOptions) (resolver.Resolver, error) {
	logx.Debugf("target: %s, callstack: %s, cc ptr: %p", target, string(debug.Stack()), cc)

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
		endpoints, err := cs.CoreV1().Endpoints(svc.Namespace).Get(context.Background(), svc.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		svc.Port = int(endpoints.Subsets[0].Ports[0].Port)
	}

	handler := kube.NewEventHandler(func(endpoints []string) {
		var addrs []resolver.Address
		for _, val := range subset(endpoints, subsetSize) {
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

	endpoints, err := cs.CoreV1().Endpoints(svc.Namespace).Get(context.Background(), svc.Name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	handler.Update(endpoints)

	resolver := &kubeResolver{
		cc:     cc,
		stopCh: make(chan struct{}),
		inf:    inf,
	}

	resolver.start()

	return resolver, nil
}

func (b *kubeBuilder) Scheme() string {
	return KubernetesScheme
}
