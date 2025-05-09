package httpx

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/validation"
	"github.com/zeromicro/go-zero/rest/internal/encoding"
	"github.com/zeromicro/go-zero/rest/internal/header"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

const (
	formKey           = "form"
	pathKey           = "path"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler = mapping.NewUnmarshaler(
		formKey,
		mapping.WithStringValues(),
		mapping.WithOpaqueKeys(),
		mapping.WithFromArray())
	pathUnmarshaler = mapping.NewUnmarshaler(
		pathKey,
		mapping.WithStringValues(),
		mapping.WithOpaqueKeys())

	// panic: sync/atomic: store of inconsistently typed value into Value
	// don't use atomic.Value to store the validator, different concrete types still panic
	validator     Validator
	validatorLock sync.RWMutex
)

// Validator defines the interface for validating the request.
type Validator interface {
	// Validate validates the request and parsed data.
	Validate(r *http.Request, data any) error
}

// Parse parses the request.
func Parse(r *http.Request, v any) error {
	kind := mapping.Deref(reflect.TypeOf(v)).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		if err := ParsePath(r, v); err != nil {
			return err
		}

		if err := ParseForm(r, v); err != nil {
			return err
		}

		if err := ParseHeaders(r, v); err != nil {
			return err
		}
	}

	if err := ParseJsonBody(r, v); err != nil {
		return err
	}

	if valid, ok := v.(validation.Validator); ok {
		return valid.Validate()
	} else if val := getValidator(); val != nil {
		return val.Validate(r, v)
	}

	return nil
}

// ParseHeaders parses the headers request.
func ParseHeaders(r *http.Request, v any) error {
	return encoding.ParseHeaders(r.Header, v)
}

// ParseForm parses the form request.
func ParseForm(r *http.Request, v any) error {
	params, err := GetFormValues(r)
	if err != nil {
		return err
	}

	return formUnmarshaler.Unmarshal(params, v)
}

// ParseHeader parses the request header and returns a map.
func ParseHeader(headerValue string) map[string]string {
	ret := make(map[string]string)
	fields := strings.Split(headerValue, separator)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}

		kv := strings.SplitN(field, "=", tokensInAttribute)
		if len(kv) != tokensInAttribute {
			continue
		}

		ret[kv[0]] = kv[1]
	}

	return ret
}

// ParseJsonBody parses the post request which contains json in body.
func ParseJsonBody(r *http.Request, v any) error {
	if withJsonBody(r) {
		reader := io.LimitReader(r.Body, maxBodyLen)
		return mapping.UnmarshalJsonReader(reader, v)
	}

	return mapping.UnmarshalJsonMap(nil, v)
}

// ParsePath parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v any) error {
	vars := pathvar.Vars(r)
	m := make(map[string]any, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

// SetValidator sets the validator.
// The validator is used to validate the request, only called in Parse,
// not in ParseHeaders, ParseForm, ParseHeader, ParseJsonBody, ParsePath.
func SetValidator(val Validator) {
	validatorLock.Lock()
	defer validatorLock.Unlock()
	validator = val
}

func getValidator() Validator {
	validatorLock.RLock()
	defer validatorLock.RUnlock()
	return validator
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(header.ContentType), header.ApplicationJson)
}
