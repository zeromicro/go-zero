package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func TestAdd(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd("bad", false)
	h.OnAdd(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
			{
				Addresses: []string{"0.0.0.3"},
			},
		},
	}, false)
	assert.ElementsMatch(t, []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, endpoints)
}

func TestDelete(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
			{
				Addresses: []string{"0.0.0.3"},
			},
		},
	}, false)
	h.OnDelete("bad")
	h.OnDelete(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
	})
	assert.ElementsMatch(t, []string{"0.0.0.3"}, endpoints)
}

func TestDeleteNamedEndpointSlice(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(newEndpointSlice("slice-a", "1", "0.0.0.1", "0.0.0.2"), false)
	h.OnAdd(newEndpointSlice("slice-b", "1", "0.0.0.3"), false)

	h.OnDelete(newEndpointSlice("slice-a", "1", "0.0.0.1", "0.0.0.2"))

	assert.ElementsMatch(t, []string{"0.0.0.3"}, endpoints)
}

func TestDeleteEndpointSliceTombstone(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(newEndpointSlice("slice-a", "1", "0.0.0.1"), false)
	h.OnAdd(newEndpointSlice("slice-b", "1", "0.0.0.2"), false)

	h.OnDelete(cache.DeletedFinalStateUnknown{
		Obj: newEndpointSlice("slice-a", "1", "0.0.0.1"),
	})

	assert.ElementsMatch(t, []string{"0.0.0.2"}, endpoints)
}

func TestUpdate(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnUpdate(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
			{
				Addresses: []string{"0.0.0.3"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, endpoints)
}

func TestUpdateNoChange(t *testing.T) {
	h := NewEventHandler(func(change []string) {
		assert.Fail(t, "should not called")
	})
	h.OnUpdate(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	})
}

func TestUpdateChangeWithDifferentVersion(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.3"},
			},
		},
	}, false)
	h.OnUpdate(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.3"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{"0.0.0.1", "0.0.0.2"}, endpoints)
}

func TestUpdateNoChangeWithDifferentVersion(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
	}, false)
	h.OnUpdate("bad", &discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
		},
	})
	h.OnUpdate(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
		},
	}, "bad")
	h.OnUpdate(&discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &discoveryv1.EndpointSlice{
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"0.0.0.1"},
			},
			{
				Addresses: []string{"0.0.0.2"},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{"0.0.0.1", "0.0.0.2"}, endpoints)
}

func TestUpdateEmptyNamedEndpointSliceKeepsOtherSlices(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(newEndpointSlice("slice-a", "1", "0.0.0.1"), false)
	h.OnAdd(newEndpointSlice("slice-b", "1", "0.0.0.2"), false)

	h.OnUpdate(newEndpointSlice("slice-a", "1", "0.0.0.1"), newEndpointSlice("slice-a", "2"))

	assert.ElementsMatch(t, []string{"0.0.0.2"}, endpoints)
}

func TestUpdateNamedEndpointSliceAggregatesAllSlices(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(newEndpointSlice("slice-a", "1", "0.0.0.1"), false)
	h.OnAdd(newEndpointSlice("slice-b", "1", "0.0.0.2"), false)

	h.OnUpdate(newEndpointSlice("slice-a", "1", "0.0.0.1"), newEndpointSlice("slice-a", "2", "0.0.0.3"))

	assert.ElementsMatch(t, []string{"0.0.0.2", "0.0.0.3"}, endpoints)
}

func newEndpointSlice(name, version string, addresses ...string) *discoveryv1.EndpointSlice {
	endpoints := make([]discoveryv1.Endpoint, 0, len(addresses))
	for _, address := range addresses {
		endpoints = append(endpoints, discoveryv1.Endpoint{
			Addresses: []string{address},
		})
	}

	return &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       "default",
			ResourceVersion: version,
		},
		Endpoints: endpoints,
	}
}
