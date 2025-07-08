package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
