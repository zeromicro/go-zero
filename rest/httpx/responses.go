package httpx

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/tal-tech/go-zero/core/logx"
)

var (
	errorHandler func(error) (int, interface{})
	lock         sync.RWMutex
)

func Error(w http.ResponseWriter, err error) {
	lock.RLock()
	handler := errorHandler
	lock.RUnlock()

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

func Ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func OkJson(w http.ResponseWriter, v interface{}) {
	WriteJson(w, http.StatusOK, v)
}

func SetErrorHandler(handler func(error) (int, interface{})) {
	lock.Lock()
	defer lock.Unlock()
	errorHandler = handler
}

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
