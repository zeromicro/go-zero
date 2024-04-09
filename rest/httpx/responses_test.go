package httpx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
		errorHandler  func(error) (int, any)
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
			errorHandler: func(err error) (int, any) {
				return http.StatusForbidden, err.Error()
			},
			expectHasBody: true,
			expectBody:    wrappedBody,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return error",
			input: body,
			errorHandler: func(err error) (int, any) {
				return http.StatusForbidden, err
			},
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return nil",
			input: body,
			errorHandler: func(err error) (int, any) {
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
				errorLock.RLock()
				prev := errorHandler
				errorLock.RUnlock()
				SetErrorHandler(test.errorHandler)
				defer func() {
					errorLock.Lock()
					errorHandler = prev
					errorLock.Unlock()
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
	t.Run("no handler", func(t *testing.T) {
		w := tracedResponseWriter{
			headers: make(map[string][]string),
		}
		msg := message{Name: "anyone"}
		OkJson(&w, msg)
		assert.Equal(t, http.StatusOK, w.code)
		assert.Equal(t, "{\"name\":\"anyone\"}", w.builder.String())
	})

	t.Run("with handler", func(t *testing.T) {
		okLock.RLock()
		prev := okHandler
		okLock.RUnlock()
		t.Cleanup(func() {
			okLock.Lock()
			okHandler = prev
			okLock.Unlock()
		})

		SetOkHandler(func(_ context.Context, v interface{}) any {
			return fmt.Sprintf("hello %s", v.(message).Name)
		})
		w := tracedResponseWriter{
			headers: make(map[string][]string),
		}
		msg := message{Name: "anyone"}
		OkJson(&w, msg)
		assert.Equal(t, http.StatusOK, w.code)
		assert.Equal(t, `"hello anyone"`, w.builder.String())
	})
}

func TestOkJsonCtx(t *testing.T) {
	t.Run("no handler", func(t *testing.T) {
		w := tracedResponseWriter{
			headers: make(map[string][]string),
		}
		msg := message{Name: "anyone"}
		OkJsonCtx(context.Background(), &w, msg)
		assert.Equal(t, http.StatusOK, w.code)
		assert.Equal(t, "{\"name\":\"anyone\"}", w.builder.String())
	})

	t.Run("with handler", func(t *testing.T) {
		okLock.RLock()
		prev := okHandler
		okLock.RUnlock()
		t.Cleanup(func() {
			okLock.Lock()
			okHandler = prev
			okLock.Unlock()
		})

		SetOkHandler(func(_ context.Context, v interface{}) any {
			return fmt.Sprintf("hello %s", v.(message).Name)
		})
		w := tracedResponseWriter{
			headers: make(map[string][]string),
		}
		msg := message{Name: "anyone"}
		OkJsonCtx(context.Background(), &w, msg)
		assert.Equal(t, http.StatusOK, w.code)
		assert.Equal(t, `"hello anyone"`, w.builder.String())
	})
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
	WriteJson(&w, http.StatusOK, map[string]any{
		"Data": complex(0, 0),
	})
	assert.Equal(t, http.StatusInternalServerError, w.code)
}

func TestStream(t *testing.T) {
	t.Run("regular case", func(t *testing.T) {
		channel := make(chan string)
		go func() {
			defer close(channel)
			for index := 0; index < 5; index++ {
				channel <- fmt.Sprintf("%d", index)
			}
		}()

		w := httptest.NewRecorder()
		Stream(context.Background(), w, func(w io.Writer) bool {
			output, ok := <-channel
			if !ok {
				return false
			}

			outputBytes := bytes.NewBufferString(output)
			_, err := w.Write(append(outputBytes.Bytes(), []byte("\n")...))
			return err == nil
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "0\n1\n2\n3\n4\n", w.Body.String())
	})

	t.Run("context done", func(t *testing.T) {
		channel := make(chan string)
		go func() {
			defer close(channel)
			for index := 0; index < 5; index++ {
				channel <- fmt.Sprintf("num: %d", index)
			}
		}()

		w := httptest.NewRecorder()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		Stream(ctx, w, func(w io.Writer) bool {
			output, ok := <-channel
			if !ok {
				return false
			}

			outputBytes := bytes.NewBufferString(output)
			_, err := w.Write(append(outputBytes.Bytes(), []byte("\n")...))
			return err == nil
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Body.String())
	})
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
		errorHandlerCtx func(context.Context, error) (int, any)
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
			errorHandlerCtx: func(ctx context.Context, err error) (int, any) {
				return http.StatusForbidden, err.Error()
			},
			expectHasBody: true,
			expectBody:    wrappedBody,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return error",
			input: body,
			errorHandlerCtx: func(ctx context.Context, err error) (int, any) {
				return http.StatusForbidden, err
			},
			expectHasBody: true,
			expectBody:    body,
			expectCode:    http.StatusForbidden,
		},
		{
			name:  "customized error handler return nil",
			input: body,
			errorHandlerCtx: func(context.Context, error) (int, any) {
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
				errorLock.RLock()
				prev := errorHandler
				errorLock.RUnlock()
				SetErrorHandlerCtx(test.errorHandlerCtx)
				defer func() {
					errorLock.Lock()
					test.errorHandlerCtx = prev
					errorLock.Unlock()
				}()
			}
			ErrorCtx(context.Background(), &w, errors.New(test.input))
			assert.Equal(t, test.expectCode, w.code)
			assert.Equal(t, test.expectHasBody, w.hasBody)
			assert.Equal(t, test.expectBody, strings.TrimSpace(w.builder.String()))
		})
	}

	// The current handler is a global event,Set default values to avoid impacting subsequent unit tests
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
	WriteJsonCtx(context.Background(), &w, http.StatusOK, map[string]any{
		"Data": complex(0, 0),
	})
	assert.Equal(t, http.StatusInternalServerError, w.code)
}
