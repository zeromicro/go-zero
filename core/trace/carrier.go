package trace

import (
	"errors"
	"net/http"
	"strings"
)

// ErrInvalidCarrier indicates an error that the carrier is invalid.
var ErrInvalidCarrier = errors.New("invalid carrier")

type (
	// Carrier interface wraps the Get and Set method.
	Carrier interface {
		Get(key string) string
		Set(key, value string)
	}

	httpCarrier http.Header
	// grpc metadata takes keys as case insensitive
	grpcCarrier map[string][]string
)

func (h httpCarrier) Get(key string) string {
	return http.Header(h).Get(key)
}

func (h httpCarrier) Set(key, val string) {
	http.Header(h).Set(key, val)
}

func (g grpcCarrier) Get(key string) string {
	if vals, ok := g[strings.ToLower(key)]; ok && len(vals) > 0 {
		return vals[0]
	}

	return ""
}

func (g grpcCarrier) Set(key, val string) {
	key = strings.ToLower(key)
	g[key] = append(g[key], val)
}
