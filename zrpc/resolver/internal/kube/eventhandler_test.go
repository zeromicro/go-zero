package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAdd(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd("bad", false)
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
				{
					IP: "0.0.0.2",
				},
				{
					IP: "0.0.0.3",
				},
			},
		},
	}}, false)
	assert.ElementsMatch(t, []string{"0.0.0.1", "0.0.0.2", "0.0.0.3"}, endpoints)
}

func TestDelete(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
				{
					IP: "0.0.0.2",
				},
				{
					IP: "0.0.0.3",
				},
			},
		},
	}}, false)
	h.OnDelete("bad")
	h.OnDelete(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
				{
					IP: "0.0.0.2",
				},
			},
		},
	}})
	assert.ElementsMatch(t, []string{"0.0.0.3"}, endpoints)
}

func TestUpdate(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
					{
						IP: "0.0.0.3",
					},
				},
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
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
				},
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
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
				{
					IP: "0.0.0.3",
				},
			},
		},
	}}, false)
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.3",
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
				},
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
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
				{
					IP: "0.0.0.2",
				},
			},
		},
	}}, false)
	h.OnUpdate("bad", &v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
			},
		},
	}})
	h.OnUpdate(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: "0.0.0.1",
				},
			},
		},
	}}, "bad")
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: "0.0.0.1",
					},
					{
						IP: "0.0.0.2",
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{"0.0.0.1", "0.0.0.2"}, endpoints)
}
