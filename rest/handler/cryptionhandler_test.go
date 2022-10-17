package handler

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/codec"
)

const (
	reqText  = "ping"
	respText = "pong"
)

var aesKey = []byte(`PdSgVkYp3s6v9y$B&E)H+MbQeThWmZq4`)

func init() {
	log.SetOutput(io.Discard)
}

func TestCryptionHandlerGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
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
		body, err := io.ReadAll(r.Body)
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
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestCryptionHandlerFlush(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(respText))
		flusher, ok := w.(http.Flusher)
		assert.True(t, ok)
		flusher.Flush()
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	expect, err := codec.EcbEncrypt(aesKey, []byte(respText))
	assert.Nil(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString(expect), recorder.Body.String())
}

func TestCryptionHandler_Hijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := newCryptionResponseWriter(resp)
	assert.NotPanics(t, func() {
		writer.Hijack()
	})

	writer = newCryptionResponseWriter(mockedHijackable{resp})
	assert.NotPanics(t, func() {
		writer.Hijack()
	})
}

func TestCryptionHandler_ContentTooLong(t *testing.T) {
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	svr := httptest.NewServer(handler)
	defer svr.Close()

	body := make([]byte, maxBytes+1)
	rand.Read(body)
	req, err := http.NewRequest(http.MethodPost, svr.URL, bytes.NewReader(body))
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
