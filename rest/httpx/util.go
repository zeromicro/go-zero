package httpx

import (
	"net/http"
)

const xForwardedFor = "X-Forwarded-For"

// GetFormValues returns the form values.
func GetFormValues(r *http.Request) (map[string]any, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return nil, err
		}
	}

	params := make(map[string]any, len(r.Form))
	for name, values := range r.Form {
		if v, ok := filterFormValues(values); ok {
			params[name] = v
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

func filterFormValues(values []string) (v any, valid bool) {
	newValues := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			newValues = append(newValues, v)
		}
	}

	switch len(newValues) {
	case 0:
		return nil, false
	case 1:
		return newValues[0], true
	default:
		return newValues, true
	}
}
