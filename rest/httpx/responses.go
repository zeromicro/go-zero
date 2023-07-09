package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/internal/errcode"
	"github.com/zeromicro/go-zero/rest/internal/header"
)

var (
	errorHandler func(context.Context, error) (int, any)
	errorLock    sync.RWMutex
	okHandler    func(context.Context, any) any
	okLock       sync.RWMutex
)

// Error writes err into w.
func Error(w http.ResponseWriter, err error, fns ...func(w http.ResponseWriter, err error)) {
	doHandleError(w, err, buildErrorHandler(context.Background()), WriteJson, fns...)
}

// ErrorCtx writes err into w.
func ErrorCtx(ctx context.Context, w http.ResponseWriter, err error,
	fns ...func(w http.ResponseWriter, err error)) {
	writeJson := func(w http.ResponseWriter, code int, v any) {
		WriteJsonCtx(ctx, w, code, v)
	}
	doHandleError(w, err, buildErrorHandler(ctx), writeJson, fns...)
}

// Ok writes HTTP 200 OK into w.
func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// OkJson writes v into w with 200 OK.
func OkJson(w http.ResponseWriter, v any) {
	okLock.RLock()
	handler := okHandler
	okLock.RUnlock()
	if handler != nil {
		v = handler(context.Background(), v)
	}
	WriteJson(w, http.StatusOK, v)
}

// OkJsonCtx writes v into w with 200 OK.
func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any) {
	okLock.RLock()
	handlerCtx := okHandler
	okLock.RUnlock()
	if handlerCtx != nil {
		v = handlerCtx(ctx, v)
	}
	WriteJsonCtx(ctx, w, http.StatusOK, v)
}

// SetErrorHandler sets the error handler, which is called on calling Error.
// Notice: SetErrorHandler and SetErrorHandlerCtx set the same error handler.
// Keeping both SetErrorHandler and SetErrorHandlerCtx is for backward compatibility.
func SetErrorHandler(handler func(error) (int, any)) {
	errorLock.Lock()
	defer errorLock.Unlock()
	errorHandler = func(_ context.Context, err error) (int, any) {
		return handler(err)
	}
}

// SetErrorHandlerCtx sets the error handler, which is called on calling Error.
// Notice: SetErrorHandler and SetErrorHandlerCtx set the same error handler.
// Keeping both SetErrorHandler and SetErrorHandlerCtx is for backward compatibility.
func SetErrorHandlerCtx(handlerCtx func(context.Context, error) (int, any)) {
	errorLock.Lock()
	defer errorLock.Unlock()
	errorHandler = handlerCtx
}

// SetOkHandler sets the response handler, which is called on calling OkJson and OkJsonCtx.
func SetOkHandler(handler func(context.Context, any) any) {
	okLock.Lock()
	defer okLock.Unlock()
	okHandler = handler
}

// WriteJson writes v as json string into w with code.
func WriteJson(w http.ResponseWriter, code int, v any) {
	if err := doWriteJson(w, code, v); err != nil {
		logx.Error(err)
	}
}

// WriteJsonCtx writes v as json string into w with code.
func WriteJsonCtx(ctx context.Context, w http.ResponseWriter, code int, v any) {
	if err := doWriteJson(w, code, v); err != nil {
		logx.WithContext(ctx).Error(err)
	}
}

func buildErrorHandler(ctx context.Context) func(error) (int, any) {
	errorLock.RLock()
	handlerCtx := errorHandler
	errorLock.RUnlock()

	var handler func(error) (int, any)
	if handlerCtx != nil {
		handler = func(err error) (int, any) {
			return handlerCtx(ctx, err)
		}
	}

	return handler
}

func doHandleError(w http.ResponseWriter, err error, handler func(error) (int, any),
	writeJson func(w http.ResponseWriter, code int, v any),
	fns ...func(w http.ResponseWriter, err error)) {
	if handler == nil {
		if len(fns) > 0 {
			for _, fn := range fns {
				fn(w, err)
			}
		} else if errcode.IsGrpcError(err) {
			// don't unwrap error and get status.Message(),
			// it hides the rpc error headers.
			http.Error(w, err.Error(), errcode.CodeFromGrpcError(err))
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		return
	}

	code, body := handler(err)
	if body == nil {
		w.WriteHeader(code)
		return
	}

	e, ok := body.(error)
	if ok {
		http.Error(w, e.Error(), code)
	} else {
		writeJson(w, code, body)
	}
}

func doWriteJson(w http.ResponseWriter, code int, v any) error {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("marshal json failed, error: %w", err)
	}

	w.Header().Set(ContentType, header.JsonContentType)
	w.WriteHeader(code)

	if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			return fmt.Errorf("write response failed, error: %w", err)
		}
	} else if n < len(bs) {
		return fmt.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}

	return nil
}
