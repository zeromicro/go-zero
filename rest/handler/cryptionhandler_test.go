package handler

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/codec"
)

const (
	reqText  = "ping"
	respText = "pong"
)

var aesKey = []byte(`PdSgVkYp3s6v9y$B&E)H+MbQeThWmZq4`)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestCryptionHandlerGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(respText))
		w.Header().Set("X-Test", "test")
		assert.Nil(t, err)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	expect, err := codec.EcbEncrypt(aesKey, []byte(respText))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "test", recorder.Header().Get("X-Test"))
	assert.Equal(t, base64.StdEncoding.EncodeToString(expect), recorder.Body.String())
}

func TestCryptionHandlerPost(t *testing.T) {
	var buf bytes.Buffer
	enc, err := codec.EcbEncrypt(aesKey, []byte(reqText))
	assert.Nil(t, err)
	buf.WriteString(base64.StdEncoding.EncodeToString(enc))

	req := httptest.NewRequest(http.MethodPost, "/any", &buf)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, reqText, string(body))

		w.Write([]byte(respText))
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	expect, err := codec.EcbEncrypt(aesKey, []byte(respText))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, base64.StdEncoding.EncodeToString(expect), recorder.Body.String())
}

func TestCryptionHandlerPostBadEncryption(t *testing.T) {
	var buf bytes.Buffer
	enc, err := codec.EcbEncrypt(aesKey, []byte(reqText))
	assert.Nil(t, err)
	buf.Write(enc)

	req := httptest.NewRequest(http.MethodPost, "/any", &buf)
	handler := CryptionHandler(aesKey)(nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCryptionHandlerWriteHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}
