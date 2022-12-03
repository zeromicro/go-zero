package httpx

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type message struct {
	Name string `json:"name"`
}

func init() {
	logx.Disable()
}

func TestError(t *testing.T) {
	const (
		body        = "foo"
		wrappedBody = `"foo"`
	)

	tests := []struct {
		name          string
		input         string
		errorHandler  func(error) (int, interface{})
		expectHasBody bool
		expectBody    string
		expectCode    int
	}{
		{
			name:          "default error handler",
			input:         body,
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusBadRequest,
		},
		{
			name:  "customized error handler return string",
			input: body,
			errorHandler: func(err error) (int, interface{}) {
				return http.StatusForbidden, err.Error()
			},
			expectHasBody: true,
			expectBody:    wrappedBody,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return error",
			input: body,
			errorHandler: func(err error) (int, interface{}) {
				return http.StatusForbidden, err
			},
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return nil",
			input: body,
			errorHandler: func(err error) (int, interface{}) {
				return http.StatusForbidden, nil
			},
			expectHasBody: false,
			expectBody:    "",
			expectCode:    http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := tracedResponseWriter{
				headers: make(map[string][]string),
			}
			if test.errorHandler != nil {
				lock.RLock()
				prev := errorHandler
				lock.RUnlock()
				SetErrorHandler(test.errorHandler)
				defer func() {
					lock.Lock()
					errorHandler = prev
					lock.Unlock()
				}()
			}
			Error(&w, errors.New(test.input))
			assert.Equal(t, test.expectCode, w.code)
			assert.Equal(t, test.expectHasBody, w.hasBody)
			assert.Equal(t, test.expectBody, strings.TrimSpace(w.builder.String()))
		})
	}
}

func TestErrorWithGrpcError(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Error(&w, status.Error(codes.Unavailable, "foo"))
	assert.Equal(t, http.StatusServiceUnavailable, w.code)
	assert.True(t, w.hasBody)
	assert.True(t, strings.Contains(w.builder.String(), "foo"))
}

func TestErrorWithHandler(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Error(&w, errors.New("foo"), func(w http.ResponseWriter, err error) {
		http.Error(w, err.Error(), 499)
	})
	assert.Equal(t, 499, w.code)
	assert.True(t, w.hasBody)
	assert.Equal(t, "foo", strings.TrimSpace(w.builder.String()))
}

func TestOk(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	Ok(&w)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestOkJson(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	OkJson(&w, msg)
	assert.Equal(t, http.StatusOK, w.code)
	assert.Equal(t, "{\"name\":\"anyone\"}", w.builder.String())
}

func TestWriteJsonTimeout(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     http.ErrHandlerTimeout,
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteJsonError(t *testing.T) {
	// only log it and ignore
	w := tracedResponseWriter{
		headers: make(map[string][]string),
		err:     errors.New("foo"),
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteJsonLessWritten(t *testing.T) {
	w := tracedResponseWriter{
		headers:     make(map[string][]string),
		lessWritten: true,
	}
	msg := message{Name: "anyone"}
	WriteJson(&w, http.StatusOK, msg)
	assert.Equal(t, http.StatusOK, w.code)
}

func TestWriteJsonMarshalFailed(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	WriteJson(&w, http.StatusOK, map[string]interface{}{
		"Data": complex(0, 0),
	})
	assert.Equal(t, http.StatusInternalServerError, w.code)
}

type tracedResponseWriter struct {
	headers     map[string][]string
	builder     strings.Builder
	hasBody     bool
	code        int
	lessWritten bool
	wroteHeader bool
	err         error
}

func (w *tracedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *tracedResponseWriter) Write(bytes []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err = w.builder.Write(bytes)
	if w.lessWritten {
		n--
	}
	w.hasBody = true

	return
}

func (w *tracedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.code = code
}

func TestErrorCtx(t *testing.T) {
	const (
		body        = "foo"
		wrappedBody = `"foo"`
	)

	tests := []struct {
		name            string
		input           string
		errorHandlerCtx func(context.Context, error) (int, interface{})
		expectHasBody   bool
		expectBody      string
		expectCode      int
	}{
		{
			name:          "default error handler",
			input:         body,
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusBadRequest,
		},
		{
			name:  "customized error handler return string",
			input: body,
			errorHandlerCtx: func(ctx context.Context, err error) (int, interface{}) {
				return http.StatusForbidden, err.Error()
			},
			expectHasBody: true,
			expectBody:    wrappedBody,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return error",
			input: body,
			errorHandlerCtx: func(ctx context.Context, err error) (int, interface{}) {
				return http.StatusForbidden, err
			},
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return nil",
			input: body,
			errorHandlerCtx: func(context.Context, error) (int, interface{}) {
				return http.StatusForbidden, nil
			},
			expectHasBody: false,
			expectBody:    "",
			expectCode:    http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := tracedResponseWriter{
				headers: make(map[string][]string),
			}
			if test.errorHandlerCtx != nil {
				lock.RLock()
				prev := errorHandlerCtx
				lock.RUnlock()
				SetErrorHandlerCtx(test.errorHandlerCtx)
				defer func() {
					lock.Lock()
					test.errorHandlerCtx = prev
					lock.Unlock()
				}()
			}
			ErrorCtx(context.Background(), &w, errors.New(test.input))
			assert.Equal(t, test.expectCode, w.code)
			assert.Equal(t, test.expectHasBody, w.hasBody)
			assert.Equal(t, test.expectBody, strings.TrimSpace(w.builder.String()))
		})
	}

	//The current handler is a global event,Set default values to avoid impacting subsequent unit tests
	SetErrorHandlerCtx(nil)
}

func TestErrorWithGrpcErrorCtx(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	ErrorCtx(context.Background(), &w, status.Error(codes.Unavailable, "foo"))
	assert.Equal(t, http.StatusServiceUnavailable, w.code)
	assert.True(t, w.hasBody)
	assert.True(t, strings.Contains(w.builder.String(), "foo"))
}

func TestErrorWithHandlerCtx(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	ErrorCtx(context.Background(), &w, errors.New("foo"), func(w http.ResponseWriter, err error) {
		http.Error(w, err.Error(), 499)
	})
	assert.Equal(t, 499, w.code)
	assert.True(t, w.hasBody)
	assert.Equal(t, "foo", strings.TrimSpace(w.builder.String()))
}

func TestWriteJsonCtxMarshalFailed(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	WriteJsonCtx(context.Background(), &w, http.StatusOK, map[string]interface{}{
		"Data": complex(0, 0),
	})
	assert.Equal(t, http.StatusInternalServerError, w.code)
}
