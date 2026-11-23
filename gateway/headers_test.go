package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyResponseHeadersDropsHopByHop(t *testing.T) {
	src := http.Header{}
	src.Set("Content-Type", "text/csv")
	src.Set("X-Request-Id", "abc")
	src.Set("Connection", "keep-alive")
	src.Set("Keep-Alive", "timeout=5")
	src.Set("Proxy-Authenticate", "Basic realm=upstream")
	src.Set("Proxy-Authorization", "Basic dXNlcjpwYXNz")
	src.Set("Upgrade", "h2c")
	src.Set("TE", "trailers")
	src.Set("Trailer", "X-Checksum")
	src.Set("Transfer-Encoding", "chunked")
	src.Set("Proxy-Connection", "keep-alive")

	dst := http.Header{}
	copyResponseHeaders(dst, src)

	assert.Equal(t, "text/csv", dst.Get("Content-Type"))
	assert.Equal(t, "abc", dst.Get("X-Request-Id"))
	for _, key := range []string{
		"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
		"Upgrade", "TE", "Trailer", "Transfer-Encoding", "Proxy-Connection",
	} {
		assert.Empty(t, dst.Values(key), key)
	}
}

func TestCopyResponseHeadersKeepsEveryValue(t *testing.T) {
	src := http.Header{}
	src.Add("Set-Cookie", "a=1")
	src.Add("Set-Cookie", "b=2")

	dst := http.Header{}
	copyResponseHeaders(dst, src)

	assert.Equal(t, []string{"a=1", "b=2"}, dst.Values("Set-Cookie"))
}

func TestHttpToHttpDropsHopByHopResponseHeaders(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("X-Request-Id", "abc")
			w.Header().Set("Keep-Alive", "timeout=5")
			w.Header().Set("Proxy-Authenticate", "Basic realm=upstream")
			w.Header().Set("Upgrade", "h2c")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("body"))
		}))
	defer upstream.Close()

	handler := new(Server).buildHttpHandler(&HttpClientConf{
		Target: upstream.Listener.Addr().String(),
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	handler.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
	assert.Equal(t, "abc", resp.Header.Get("X-Request-Id"))
	// The client keeps its connection to the gateway, not to the upstream.
	for _, key := range []string{"Keep-Alive", "Proxy-Authenticate", "Upgrade"} {
		assert.Empty(t, resp.Header.Values(key), key)
	}
	assert.Equal(t, "body", w.Body.String())
}
