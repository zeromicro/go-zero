package encoding

import (
	"net/http"
	"net/textproto"

	"github.com/zeromicro/go-zero/core/mapping"
)

const headerKey = "header"

var headerUnmarshaler = mapping.NewUnmarshaler(headerKey, mapping.WithStringValues(),
	mapping.WithCanonicalKeyFunc(textproto.CanonicalMIMEHeaderKey))

// ParseHeaders parses the headers request.
func ParseHeaders(header http.Header, v interface{}) error {
	m := map[string]interface{}{}
	for k, v := range header {
		if len(v) == 1 {
			m[k] = v[0]
		} else {
			m[k] = v
		}
	}

	return headerUnmarshaler.Unmarshal(m, v)
}
