package httpx

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"zero/core/httprouter"
	"zero/core/mapping"
)

const (
	multipartFormData = "multipart/form-data"
	xForwardFor       = "X-Forward-For"
	formKey           = "form"
	pathKey           = "path"
	emptyJson         = "{}"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	ErrBodylessRequest = errors.New("not a POST|PUT|PATCH request")

	formUnmarshaler = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
)

// Returns the peer address, supports X-Forward-For
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardFor)
	if len(v) > 0 {
		return v
	}
	return r.RemoteAddr
}

func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// Parses the form request.
func ParseForm(r *http.Request, v interface{}) error {
	if strings.Index(r.Header.Get(ContentType), multipartFormData) != -1 {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			return err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	return formUnmarshaler.Unmarshal(params, v)
}

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

// Parses the post request which contains json in body.
func ParseJsonBody(r *http.Request, v interface{}) error {
	var reader io.Reader

	if withJsonBody(r) {
		reader = io.LimitReader(r.Body, maxBodyLen)
	} else {
		reader = strings.NewReader(emptyJson)
	}

	return mapping.UnmarshalJsonReader(reader, v)
}

// Parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := httprouter.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Index(r.Header.Get(ContentType), ApplicationJson) != -1
}
