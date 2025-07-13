package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestNewServerless(t *testing.T) {
	logtest.Discard(t)

	const configYaml = `
Name: foo
Host: localhost
Port: 0
`
	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	svr, err := NewServer(cnf)
	assert.NoError(t, err)

	svr.AddRoute(Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		},
	})

	serverless, err := NewServerless(svr)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	serverless.Serve(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello World", w.Body.String())
}

func TestNewServerlessWithError(t *testing.T) {
	logtest.Discard(t)

	const configYaml = `
Name: foo
Host: localhost
Port: 0
`
	var cnf RestConf
	assert.Nil(t, conf.LoadFromYamlBytes([]byte(configYaml), &cnf))

	svr, err := NewServer(cnf)
	assert.NoError(t, err)

	svr.AddRoute(Route{
		Method:  http.MethodGet,
		Path:    "notstartwith/",
		Handler: nil,
	})

	_, err = NewServerless(svr)
	assert.Error(t, err)
}
