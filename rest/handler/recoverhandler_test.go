package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestWithPanic(t *testing.T) {
	t.Run("no recoverCallback", func(t *testing.T) {
		handler := RecoverHandler(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("whatever")
		}))

		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("recoverCallback", func(t *testing.T) {
		handler := RecoverHandler(func(w http.ResponseWriter, r *http.Request, recoverRes interface{}) {
			w.WriteHeader(http.StatusBadRequest)
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("whatever")
		}))

		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("recoverCallback panic", func(t *testing.T) {
		handler := RecoverHandler(func(w http.ResponseWriter, r *http.Request, recoverRes interface{}) {
			panic("any")
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("whatever")
		}))

		req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestWithoutPanic(t *testing.T) {
	handler := RecoverHandler(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
