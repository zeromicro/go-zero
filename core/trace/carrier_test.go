package trace

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/3Rivers/go-zero/core/stringx"
	"github.com/stretchr/testify/assert"
)

func TestHttpCarrier(t *testing.T) {
	tests := []map[string]string{
		{},
		{
			"first":  "a",
			"second": "b",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			carrier := httpCarrier(req.Header)
			for k, v := range test {
				carrier.Set(k, v)
			}
			for k, v := range test {
				assert.Equal(t, v, carrier.Get(k))
			}
			assert.Equal(t, "", carrier.Get("none"))
		})
	}
}

func TestGrpcCarrier(t *testing.T) {
	tests := []map[string]string{
		{},
		{
			"first":  "a",
			"second": "b",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			m := make(map[string][]string)
			carrier := grpcCarrier(m)
			for k, v := range test {
				carrier.Set(k, v)
			}
			for k, v := range test {
				assert.Equal(t, v, carrier.Get(k))
			}
			assert.Equal(t, "", carrier.Get("none"))
		})
	}
}
