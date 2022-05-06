package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
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
	in.Informer().AddEventHandler(handler)
	threading.GoSafe(func() {
		inf.Start(proc.Done())
	})

	endpoints, err := cs.CoreV1().Endpoints(svc.Namespace).Get(context.Background(), svc.Name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	handler.Update(endpoints)

	return &nopResolver{cc: cc}, nil
}

func (b *kubeBuilder) Scheme() string {
	return KubernetesScheme
}
