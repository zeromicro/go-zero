package kube

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect Service
		hasErr bool
	}{
		{
			name:  "normal case",
			input: "k8s://ns1/my-svc:8080",
			expect: Service{
				Namespace: "ns1",
				Name:      "my-svc",
				Port:      8080,
			},
		},
		{
			name:  "normal case",
			input: "k8s:///my-svc:8080",
			expect: Service{
				Namespace: defaultNamespace,
				Name:      "my-svc",
				Port:      8080,
			},
		},
		{
			name:   "no port",
			input:  "k8s://ns1/my-svc:",
			hasErr: true,
		},
		{
			name:  "no port, no colon",
			input: "k8s://ns1/my-svc",
			expect: Service{
				Namespace: "ns1",
				Name:      "my-svc",
			},
		},
		{
			name:   "bad port",
			input:  "k8s://ns1/my-svc:800a",
			hasErr: true,
		},
		{
			name:   "bad endpoint",
			input:  "k8s://ns1:800/:",
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := url.Parse(test.input)
			if assert.NoError(t, err) {
				svc, err := ParseTarget(resolver.Target{URL: *uri})
				if test.hasErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
					assert.Equal(t, test.expect, svc)
				}
			}
		})
	}
}
