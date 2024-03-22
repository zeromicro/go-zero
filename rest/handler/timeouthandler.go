package handler

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	statusClientClosedRequest = 499
	reason                    = "Request Timeout"
	headerUpgrade             = "Upgrade"
	valueWebsocket            = "websocket"
	headerAccept              = "Accept"
	valueSSE                  = "text/event-stream"
)

// TimeoutHandler returns the handler with given timeout.
// If client closed request, code 499 will be logged.
// Notice: even if canceled in server side, 499 will be logged as well.
func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration <= 0 {
			return next
		}

		return &timeoutHandler{
			handler: next,
			dt:      duration,
		}
	}
}

// timeoutHandler is the handler that controls the request timeout.
// Why we implement it on our own, because the stdlib implementation
// treats the ClientClosedRequest as http.StatusServiceUnavailable.
// And we write the codes in logs as code 499, which is defined by nginx.
type timeoutHandler struct {
	handler http.Handler
	dt      time.Duration
}

func (h *timeoutHandler) errorBody() string {
	return reason
}

func (h *timeoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(headerUpgrade) == valueWebsocket ||
		// Server-Sent Event ignore timeout.
		r.Header.Get(headerAccept) == valueSSE {
		h.handler.ServeHTTP(w, r)
		return
	}
	ctx, cancelCtx := context.WithCancel(r.Context())
	r = r.WithContext(ctx)
	// creat a timer
	t := time.NewTimer(h.dt)
	done := make(chan struct{})
	tw := &timeoutWriter{
		w: w,
		resetTimer: func() {
			t.Reset(h.dt)
		},
	}
	defer t.Stop()
	defer cancelCtx()
	panicChan := make(chan any, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		h.handler.ServeHTTP(tw, r)
		close(done)
	}()
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		tw.mu.Lock()
		defer tw.mu.Unlock()
		tw.Flush()
	case <-t.C:
		tw.mu.Lock()
		defer tw.mu.Unlock()
		// there isn't any user-defined middleware before TimoutHandler,
		// so we can guarantee that cancelation in biz related code won't come here.
		httpx.ErrorCtx(r.Context(), w, r.Context().Err(), func(w http.ResponseWriter, err error) {
			if errors.Is(err, context.Canceled) {
				tw.WriteHeader(statusClientClosedRequest)
			} else {
				tw.WriteHeader(http.StatusServiceUnavailable)
			}
			_, _ = io.WriteString(w, h.errorBody())
		})
		tw.timedOut = true
	}
}

type timeoutWriter struct {
	w http.ResponseWriter

	resetTimer  func()
	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
}

var _ http.Pusher = (*timeoutWriter)(nil)

// Flush implements the Flusher interface.
func (tw *timeoutWriter) Flush() {
	flusher, ok := tw.w.(http.Flusher)
	if !ok {
		return
	}
	tw.resetTimer()
	flusher.Flush()
}

// Header returns the underline temporary http.Header.
func (tw *timeoutWriter) Header() http.Header {
	return tw.w.Header()
}

// Hijack implements the Hijacker interface.
func (tw *timeoutWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := tw.w.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

// Push implements the Pusher interface.
func (tw *timeoutWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := tw.w.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}

	return http.ErrNotSupported
}

// Write writes the data to the connection as part of an HTTP reply.
// Timeout and multiple header written are guarded.
func (tw *timeoutWriter) Write(p []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	return tw.w.Write(p)
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if !tw.wroteHeader {
		tw.w.WriteHeader(code)
		tw.wroteHeader = true
	}
}
