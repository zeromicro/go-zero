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
	errorHandler    func(error) (int, any)
	errorHandlerCtx func(context.Context, error) (int, any)
	lock            sync.RWMutex
)

// Error writes err into w.
func Error(w http.ResponseWriter, err error, fns ...func(w http.ResponseWriter, err error)) {
	lock.RLock()
	handler := errorHandler
	lock.RUnlock()

	doHandleError(w, err, handler, WriteJson, fns...)
}

// ErrorCtx writes err into w.
func ErrorCtx(ctx context.Context, w http.ResponseWriter, err error,
	fns ...func(w http.ResponseWriter, err error)) {
	lock.RLock()
	handlerCtx := errorHandlerCtx
	lock.RUnlock()

	var handler func(error) (int, any)
	if handlerCtx != nil {
		handler = func(err error) (int, any) {
			return handlerCtx(ctx, err)
		}
	}
	writeJson := func(w http.ResponseWriter, code int, v any) {
		WriteJsonCtx(ctx, w, code, v)
	}
	doHandleError(w, err, handler, writeJson, fns...)
}

// Ok writes HTTP 200 OK into w.
func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// OkJson writes v into w with 200 OK.
func OkJson(w http.ResponseWriter, v any) {
	WriteJson(w, http.StatusOK, v)
}

// OkJsonCtx writes v into w with 200 OK.
func OkJsonCtx(ctx context.Context, w http.ResponseWriter, v any) {
	WriteJsonCtx(ctx, w, http.StatusOK, v)
}

// SetErrorHandler sets the error handler, which is called on calling Error.
func SetErrorHandler(handler func(error) (int, any)) {
	lock.Lock()
	defer lock.Unlock()
	errorHandler = handler
}

// SetErrorHandlerCtx sets the error handler, which is called on calling Error.
func SetErrorHandlerCtx(handlerCtx func(context.Context, error) (int, any)) {
	lock.Lock()
	defer lock.Unlock()
	errorHandlerCtx = handlerCtx
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
