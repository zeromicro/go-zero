package gateway

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHeadersNoValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("a", "b")
	assert.Nil(t, buildHeaders(req.Header))
}

func TestBuildHeadersWithValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("grpc-metadata-a", "b")
	req.Header.Add("grpc-metadata-b", "b")
	assert.EqualValues(t, []string{"gateway-A:b", "gateway-B:b"}, buildHeaders(req.Header))
}
