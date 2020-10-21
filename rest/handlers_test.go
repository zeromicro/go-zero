package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsHandler(t *testing.T) {
	w := httptest.NewRecorder()
	handler := CorsHandler()
	handler.ServeHTTP(w, nil)
	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	assert.Equal(t, allOrigin, w.Header().Get(allowOrigin))
}

func TestCorsHandlerWithOrigins(t *testing.T) {
	origins := []string{"local", "remote"}
	w := httptest.NewRecorder()
	handler := CorsHandler(origins...)
	handler.ServeHTTP(w, nil)
	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	assert.Equal(t, strings.Join(origins, separator), w.Header().Get(allowOrigin))
}
