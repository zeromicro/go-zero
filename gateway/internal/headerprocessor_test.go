package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHeadersNoValue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Add("a", "b")
	assert.Nil(t, ProcessHeaders(req.Header))
}

func TestBuildHeadersWithValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Add("grpc-metadata-a", "b")
	req.Header.Add("grpc-metadata-b", "b")
	assert.ElementsMatch(t, []string{"gateway-A:b", "gateway-B:b"}, ProcessHeaders(req.Header))
}
