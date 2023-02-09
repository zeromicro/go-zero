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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := url.Parse(test.input)
			assert.Nil(t, err)
			svc, err := ParseTarget(resolver.Target{URL: *uri})
			if test.hasErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expect, svc)
			}
		})
	}
}

func TestParseNonBlock(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expect   bool
		emptySvc bool
	}{
		{
			name:     "no block",
			input:    "k8s://ns1/my-svc:8080",
			expect:   false,
			emptySvc: false,
		},
		{
			name:     "block with true",
			input:    "k8s://ns1:8080?nonBlock=true",
			expect:   true,
			emptySvc: false,
		},
		{
			name:     "block with false",
			input:    "k8s://ns1:8080?nonBlock=false",
			expect:   false,
			emptySvc: false,
		},
		{
			name:     "block with error param",
			input:    "k8s://ns1:8080?nonBlock=abcd",
			expect:   false,
			emptySvc: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			uri, err := url.Parse(test.input)
			assert.Nil(t, err)

			target := resolver.Target{URL: *uri}
			svc, err := ParseTarget(target)
			if test.emptySvc {
				assert.Equal(t, emptyService, svc)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, svc.NonBlock, test.expect)
			}
		})
	}

}
