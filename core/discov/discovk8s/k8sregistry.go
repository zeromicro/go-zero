package discovk8s

import (
	"github.com/tal-tech/go-zero/core/logx"
	"os"
	"path/filepath"
	"sync"
)

type K8sRegistry struct {
	Registry

	k8sEpController *EndpointController

	services map[string][]*ServiceInstance

	lock sync.Mutex
}

func NewK8sRegistry(config *K8sConfig) *K8sRegistry {
	kubeconfig := getKubeconfig(config)

	k8sClient, err := NewK8sClient(kubeconfig)

	if err != nil {
		logx.Errorf("new k8s client error: %v", err)
		os.Exit(1)
	}

	k8sEpController := NewEndpointController(k8sClient)

	registry := K8sRegistry{
		k8sEpController: k8sEpController,
		services:        make(map[string][]*ServiceInstance),
	}

	return &registry
}

func (r *K8sRegistry) NewSubscriber(service *Service) Subscriber {

	serviceFullName := service.EpName()

	k8sSubscriber := K8sSubscriber{
		k8sRegistry: r,
		instance:    service,
	}

	r.k8sEpController.AddOnUpdateFunc(serviceFullName, func(addresses []*ServiceInstance) {
		r.lock.Lock()
		r.services[serviceFullName] = addresses
		r.lock.Unlock()

		k8sSubscriber.OnUpdate()
	})

	return &k8sSubscriber
}

func (r *K8sRegistry) GetServices(service *Service) []*ServiceInstance {
	r.lock.Lock()
	defer r.lock.Unlock()

	if value, ok := r.services[service.EpName()]; ok {
		var ret []*ServiceInstance

		for _, item := range value {
			if item.Port == service.Port {
				ret = append(ret, item)
			}
		}
		return ret
	} else {
		fullService, _ := r.k8sEpController.GetEndpoints(service.Name, service.Namespace)
		var ret []*ServiceInstance

		for _, item := range fullService {
			if item.Port == service.Port {
				ret = append(ret, item)
			}
		}
		return ret
	}

}

func getKubeconfig(config *K8sConfig) string {

	if config != nil && config.KubeconfigFile != "" {
		return config.KubeconfigFile
	}
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	_, err := os.Lstat(kubeconfig)

	if err == nil || os.IsExist(err) {
		return kubeconfig
	}
	return ""
}
