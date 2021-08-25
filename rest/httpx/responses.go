package httpx

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/tal-tech/go-zero/core/logx"
)

var (
	successHandler func(interface{}) interface{}
	errorHandler   func(error) (int, interface{})
	errLock        sync.RWMutex
	sucLock        sync.RWMutex
)

// Error writes err into w.
func Error(w http.ResponseWriter, err error) {
	errLock.RLock()
	handler := errorHandler
	errLock.RUnlock()

	if handler == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	code, body := errorHandler(err)
	e, ok := body.(error)
	if ok {
		http.Error(w, e.Error(), code)
	} else {
		WriteJson(w, code, body)
	}
}

func Success(w http.ResponseWriter, data interface{}) {
	sucLock.RLock()
	handler := successHandler
	sucLock.RUnlock()

	if handler == nil {
		OkJson(w, data)
		return
	}

	body := successHandler(data)
	OkJson(w, body)
}

// Ok writes HTTP 200 OK into w.
func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

// OkJson writes v into w with 200 OK.
func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

// SetErrorHandler sets the error handler, which is called on calling Error.
func SetErrorHandler(handler func(error) (int, interface{})) {
	errLock.Lock()
	defer errLock.Unlock()
	errorHandler = handler
}

func SetSuccessHandler(handler func(interface{}) interface{}) {
	sucLock.Lock()
	defer sucLock.Unlock()
	successHandler = handler
}

// WriteJson writes v as json string into w with code.
func WriteJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(code)

	if bs, err := json.Marshal(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if n, err := w.Write(bs); err != nil {
		// http.ErrHandlerTimeout has been handled by http.TimeoutHandler,
		// so it's ignored here.
		if err != http.ErrHandlerTimeout {
			logx.Errorf("write response failed, error: %s", err)
		}
	} else if n < len(bs) {
		logx.Errorf("actual bytes: %d, written bytes: %d", len(bs), n)
	}
}
