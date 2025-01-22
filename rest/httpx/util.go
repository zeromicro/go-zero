package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	xForwardedFor = "X-Forwarded-For"
	arraySuffix   = "[]"
	// most servers and clients have a limit of 8192 bytes (8 KB)
	// one parameter at least take 4 chars, for example `?a=b&c=d`
	maxFormParamCount = 2048
)

// GetFormValues returns the form values supporting three array notation formats:
//  1. Standard notation: /api?names=alice&names=bob
//  2. Comma notation: /api?names=alice,bob
//  3. Bracket notation: /api?names[]=alice&names[]=bob
func GetFormValues(r *http.Request) (map[string]any, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return nil, err
		}
	}

	var n int
	params := make(map[string]any, len(r.Form))
	for name, values := range r.Form {
		filtered := make([]string, 0, len(values))
		for _, v := range values {
			// ignore empty values, especially for optional int parameters
			// e.g. /api?ids=
			// e.g. /api
			// type Req struct {
			//	IDs []int `form:"ids,optional"`
			// }
			if len(v) == 0 {
				continue
			}

			if n < maxFormParamCount {
				filtered = append(filtered, v)
				n++
			} else {
				return nil, fmt.Errorf("too many form values, error: %s", r.Form.Encode())
			}
		}

		if len(filtered) > 0 {
			if strings.HasSuffix(name, arraySuffix) {
				name = name[:len(name)-2]
			}
			params[name] = filtered
		}
	}

	return params, nil
}

// GetRemoteAddr returns the peer address, supports X-Forward-For.
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardedFor)
	if len(v) > 0 {
		return v
	}

	return r.RemoteAddr
}
