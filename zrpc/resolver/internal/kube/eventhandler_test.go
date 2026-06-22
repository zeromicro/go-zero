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
	falseVal := false
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
			{
				Addresses: []string{"0.0.0.4"},
				Conditions: discoveryv1.EndpointConditions{
					Ready: &falseVal,
				},
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
	falseVal := false
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
			{
				Addresses: []string{"0.0.0.4"},
				Conditions: discoveryv1.EndpointConditions{
					Ready: &falseVal,
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

func TestIsValidEndpoint(t *testing.T) {
	trueVal := true
	falseVal := false

	tests := []struct {
		name     string
		point    discoveryv1.Endpoint
		expected bool
	}{
		{
			name:     "all nil conditions",
			point:    discoveryv1.Endpoint{},
			expected: true,
		},
		{
			name: "ready true",
			point: discoveryv1.Endpoint{
				Conditions: discoveryv1.EndpointConditions{
					Ready: &trueVal,
				},
			},
			expected: true,
		},
		{
			name: "ready false",
			point: discoveryv1.Endpoint{
				Conditions: discoveryv1.EndpointConditions{
					Ready: &falseVal,
				},
			},
			expected: false,
		},
		{
			name: "terminating true",
			point: discoveryv1.Endpoint{
				Conditions: discoveryv1.EndpointConditions{
					Terminating: &trueVal,
				},
			},
			expected: false,
		},
		{
			name: "terminating false",
			point: discoveryv1.Endpoint{
				Conditions: discoveryv1.EndpointConditions{
					Terminating: &falseVal,
				},
			},
			expected: true,
		},
		{
			name: "serving false",
			point: discoveryv1.Endpoint{
				Conditions: discoveryv1.EndpointConditions{
					Serving: &falseVal,
				},
			},
			expected: false,
		},
		{
			name: "serving true",
			point: discoveryv1.Endpoint{
				Conditions: discoveryv1.EndpointConditions{
					Serving: &trueVal,
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, isValidEndpoint(test.point))
		})
	}
}
