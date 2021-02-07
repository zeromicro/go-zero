package discovk8s

import (
	"github.com/tal-tech/go-zero/core/logx"
	v1 "k8s.io/api/core/v1"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"sort"
	"sync"
)

type EndpointController struct {
	endpointsInformer cache.SharedIndexInformer
	endpointsLister   listerv1.EndpointsLister
	updatefuncs       map[string][]OnUpdateFunc
	lock              sync.Mutex
}

type OnUpdateFunc func(addresses []*ServiceInstance)

func NewEndpointController(k8sClient *K8sClient) *EndpointController {
	informer := k8sClient.InformerFactory.Core().V1().Endpoints().Informer()
	lister := k8sClient.InformerFactory.Core().V1().Endpoints().Lister()

	c := EndpointController{
		endpointsInformer: informer,
		endpointsLister:   lister,
		updatefuncs:       make(map[string][]OnUpdateFunc),
	}

	informer.AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.endpointAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.endpointUpdate,
			// Called on resource deletion.
			DeleteFunc: c.endpointDelete,
		},
	)

	k8sClient.InformerFactory.Start(k8sClient.stop)

	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(k8sClient.stop, informer.HasSynced) {
		logx.Error("Failed to sync")
		return nil
	}

	return &c
}

func (e *EndpointController) AddOnUpdateFunc(key string, l OnUpdateFunc) {
	e.lock.Lock()
	defer e.lock.Unlock()

	if _, ok := e.updatefuncs[key]; !ok {
		e.updatefuncs[key] = make([]OnUpdateFunc, 1)
	}

	e.updatefuncs[key] = append(e.updatefuncs[key], l)

}

func (e *EndpointController) GetEndpoints(name string, namespace string) ([]*ServiceInstance, error) {
	ep, err := e.endpointsLister.Endpoints(namespace).Get(name)
	if err != nil {
		logx.Errorf("get endpoints error, %v", err)
		return nil, err
	}

	return getReadyAddress(ep), nil
}

func (e *EndpointController) endpointAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)

	epName := buildEpName(endpoints.Name, endpoints.Namespace)
	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}

	e.lock.Lock()
	copyListeners := append([]OnUpdateFunc(nil), e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, l := range copyListeners {
		l(getReadyAddress(endpoints))
	}
}

func (e *EndpointController) endpointUpdate(old, new interface{}) {
	oldEndpoint := old.(*v1.Endpoints)
	newEndpoint := new.(*v1.Endpoints)

	epName := buildEpName(oldEndpoint.Name, newEndpoint.Namespace)

	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}

	oldAddress := getReadyAddress(oldEndpoint)
	newAddress := getReadyAddress(newEndpoint)

	if reflect.DeepEqual(oldAddress, newAddress) {
		return
	}

	e.lock.Lock()
	var copyListeners []OnUpdateFunc
	copyListeners = append(copyListeners, e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, l := range copyListeners {
		if l != nil {
			l(newAddress)
		}
	}

}

func (e *EndpointController) endpointDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)

	epName := buildEpName(endpoints.Name, endpoints.Namespace)
	if _, ok := e.updatefuncs[epName]; !ok {
		return
	}
	e.lock.Lock()
	copyListeners := append([]OnUpdateFunc(nil), e.updatefuncs[epName]...)
	e.lock.Unlock()

	for _, l := range copyListeners {
		l(nil)
	}
}

func buildEpName(name string, namespace string) string {
	return name + "." + namespace
}

func getReadyAddress(endpoints *v1.Endpoints) []*ServiceInstance {
	var readyAddesses []*ServiceInstance

	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				si := ServiceInstance{
					Ip:   address.IP,
					Port: port.Port,
				}
				readyAddesses = append(readyAddesses, &si)
			}
		}
	}

	sort.Sort(ServiceInstanceSlice(readyAddesses))

	return readyAddesses
}
