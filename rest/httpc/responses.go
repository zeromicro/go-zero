package httpc

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/rest/internal/encoding"
	"github.com/zeromicro/go-zero/rest/internal/header"
)

// Parse parses the response.
func Parse(resp *http.Response, val any) error {
	if err := ParseHeaders(resp, val); err != nil {
		return err
	}

	return ParseJsonBody(resp, val)
}

// ParseHeaders parses the response headers.
func ParseHeaders(resp *http.Response, val any) error {
	return encoding.ParseHeaders(resp.Header, val)
}

// ParseJsonBody parses the response body, which should be in json content type.
func ParseJsonBody(resp *http.Response, val any) error {
	defer resp.Body.Close()

	if isContentTypeJson(resp) {
		if resp.ContentLength > 0 {
			return mapping.UnmarshalJsonReader(resp.Body, val)
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, resp.Body); err != nil {
			return err
		}

		if buf.Len() > 0 {
			return mapping.UnmarshalJsonReader(&buf, val)
		}
	}

	return mapping.UnmarshalJsonMap(nil, val)
}

func isContentTypeJson(r *http.Response) bool {
	return strings.Contains(r.Header.Get(header.ContentType), header.ApplicationJson)
}
