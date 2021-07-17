package httpx

import (
	"io"
	"net/http"
	"net/textproto"
	"strings"

	"github.com/tal-tech/go-zero/core/mapping"
	"github.com/tal-tech/go-zero/rest/internal/context"
)

const (
	formKey           = "form"
	pathKey           = "path"
	headerKey         = "header"
	emptyJson         = "{}"
	maxMemory         = 32 << 20 // 32MB
	maxBodyLen        = 8 << 20  // 8MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler   = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler   = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
	headerUnmarshaler = mapping.NewUnmarshaler(headerKey, mapping.WithStringValues())
)

// Parse parses the request.
func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	if err := ParseHeaders(r, v); err != nil {
		return err
	}

	return ParseJsonBody(r, v)
}

// ParseHeaders parses the headers request.
func ParseHeaders(r *http.Request, v interface{}) error {
	m := make(mapping.HeaderMapValuer, len(r.Header))
	for k, v := range r.Header {
		m[textproto.CanonicalMIMEHeaderKey(k)] = v
	}

	return headerUnmarshaler.UnmarshalValuer(m, v)
}

// ParseForm parses the form request.
func ParseForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}

	params := make(mapping.FormMapValuer, len(r.Form))
	for name, value := range r.Form {
		if len(value) > 0 {
			params[name] = value
		}
	}

	return formUnmarshaler.UnmarshalValuer(params, v)
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
func ParseJsonBody(r *http.Request, v interface{}) error {
	var reader io.Reader
	if withJsonBody(r) {
		reader = io.LimitReader(r.Body, maxBodyLen)
	} else {
		reader = strings.NewReader(emptyJson)
	}

	return mapping.UnmarshalJsonReader(reader, v)
}

// ParsePath parses the symbols reside in url path.
// Like http://localhost/bag/:name
func ParsePath(r *http.Request, v interface{}) error {
	vars := context.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, v := range vars {
		m[k] = v
	}

	return pathUnmarshaler.Unmarshal(m, v)
}

func withJsonBody(r *http.Request) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(ContentType), ApplicationJson)
}
