package internal

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartHttp(t *testing.T) {
	svr := httptest.NewUnstartedServer(http.NotFoundHandler())
	fields := strings.Split(svr.Listener.Addr().String(), ":")
	port, err := strconv.Atoi(fields[1])
	assert.Nil(t, err)
	err = StartHttp(fields[0], port, http.NotFoundHandler(), func(svr *http.Server) {
		svr.IdleTimeout = 0
	})
	assert.NotNil(t, err)
}

func TestStartHttps(t *testing.T) {
	svr := httptest.NewTLSServer(http.NotFoundHandler())
	fields := strings.Split(svr.Listener.Addr().String(), ":")
	port, err := strconv.Atoi(fields[1])
	assert.Nil(t, err)
	err = StartHttps(fields[0], port, "", "", http.NotFoundHandler(), func(svr *http.Server) {
		svr.IdleTimeout = 0
	})
	assert.NotNil(t, err)
}
