package discovk8s

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"strconv"
	"time"

	"path/filepath"

	"os"
	"testing"
)

// EndpointLoggingController logs the name and namespace of pods that are added,
// deleted, or updated
type EndpointLoggingController struct {
	informerFactory  informers.SharedInformerFactory
	endpointInformer coreinformers.EndpointsInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *EndpointLoggingController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.endpointInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *EndpointLoggingController) endpointAdd(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	fmt.Printf("Endpoint CREATED: %s/%s\n", endpoints.Namespace, endpoints.Name)
}

func (c *EndpointLoggingController) endpointUpdate(old, new interface{}) {
	oldEndpoint := old.(*v1.Endpoints)
	newEndpoint := new.(*v1.Endpoints)

	var oldAddesses []string
	var newAddesses []string

	for _, subset := range oldEndpoint.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				oldAddesses = append(oldAddesses, address.IP+":"+strconv.Itoa(int(port.Port)))
			}
		}
	}

	for _, subset := range newEndpoint.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				newAddesses = append(newAddesses, address.IP+":"+strconv.Itoa(int(port.Port)))
			}
		}
	}

	if !reflect.DeepEqual(oldAddesses, newAddesses) {
		fmt.Printf(
			"Endpoint UPDATED. %v/%v\n",
			oldAddesses, newAddesses,
		)
	}

}

func (c *EndpointLoggingController) endpointDelete(obj interface{}) {
	endpoints := obj.(*v1.Endpoints)
	fmt.Printf("POD DELETED: %s/%s\n", endpoints.Namespace, endpoints.Name)
}

// NewEndpointLoggingController creates a EndpointLoggingController
func NewEndpointLoggingController(informerFactory informers.SharedInformerFactory) *EndpointLoggingController {
	endpointsInformer := informerFactory.Core().V1().Endpoints()

	c := &EndpointLoggingController{
		informerFactory:  informerFactory,
		endpointInformer: endpointsInformer,
	}
	endpointsInformer.Informer().AddEventHandler(
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
	return c
}

func TestK8sClient_GetPodsExample(t *testing.T) {
	k8sClient, err := NewK8sClient(localKubeconfig())
	assert.Nil(t, err)

	ctx := context.Background()
	pods, err := k8sClient.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})

	assert.Nil(t, err)

	for _, pod := range pods.Items {
		fmt.Println(pod.Namespace, pod.Name, pod.Status.Phase)
	}
}

func TestK8sClient_WatchEndpointsExample(t *testing.T) {
	k8sClient, err := NewK8sClient(localKubeconfig())
	assert.Nil(t, err)

	factory := informers.NewSharedInformerFactory(k8sClient.Clientset, time.Hour*24)
	controller := NewEndpointLoggingController(factory)
	stop := make(chan struct{})
	defer close(stop)
	controller.Run(stop)

	select {}
}

func localKubeconfig() string {
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

func TestNewK8sClient(t *testing.T) {
	k8sClient, err := NewK8sClient(localKubeconfig())

	assert.NotNil(t, k8sClient)
	assert.Nil(t, err)
}
