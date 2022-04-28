package httpc

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/rest/internal/encoding"
	"github.com/zeromicro/go-zero/rest/internal/header"
)

// Parse parses the response.
func Parse(resp *http.Response, val interface{}) error {
	if err := ParseHeaders(resp, val); err != nil {
		return err
	}

	return ParseJsonBody(resp, val)
}

// ParseHeaders parses the rsponse headers.
func ParseHeaders(resp *http.Response, val interface{}) error {
	return encoding.ParseHeaders(resp.Header, val)
}

// ParseJsonBody parses the rsponse body, which should be in json content type.
func ParseJsonBody(resp *http.Response, val interface{}) error {
	if withJsonBody(resp) {
		return mapping.UnmarshalJsonReader(resp.Body, val)
	}

	return mapping.UnmarshalJsonMap(nil, val)
}

func withJsonBody(r *http.Response) bool {
	return r.ContentLength > 0 && strings.Contains(r.Header.Get(header.ContentType), header.ApplicationJson)
}
