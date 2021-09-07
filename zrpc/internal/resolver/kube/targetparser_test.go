package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name   string
		input  resolver.Target
		expect Service
		hasErr bool
	}{
		{
			name: "normal case",
			input: resolver.Target{
				Scheme:    "k8s",
				Authority: "ns1",
				Endpoint:  "my-svc:8080",
			},
			expect: Service{
				Namespace: "ns1",
				Name:      "my-svc",
				Port:      8080,
			},
		},
		{
			name: "normal case",
			input: resolver.Target{
				Scheme:    "k8s",
				Authority: "",
				Endpoint:  "my-svc:8080",
			},
			expect: Service{
				Namespace: defaultNamespace,
				Name:      "my-svc",
				Port:      8080,
			},
		},
		{
			name: "no port",
			input: resolver.Target{
				Scheme:    "k8s",
				Authority: "ns1",
				Endpoint:  "my-svc:",
			},
			hasErr: true,
		},
		{
			name: "no port, no colon",
			input: resolver.Target{
				Scheme:    "k8s",
				Authority: "ns1",
				Endpoint:  "my-svc",
			},
			hasErr: true,
		},
		{
			name: "bad port",
			input: resolver.Target{
				Scheme:    "k8s",
				Authority: "ns1",
				Endpoint:  "my-svc:800a",
			},
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			svc, err := ParseTarget(test.input)
			if test.hasErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expect, svc)
			}
		})
	}
}
