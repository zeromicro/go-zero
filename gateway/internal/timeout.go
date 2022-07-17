package internal

import (
	"net/http"
	"time"
)

const grpcTimeoutHeader = "Grpc-Timeout"

// GetTimeout returns the timeout from the header, if not set, returns the default timeout.
func GetTimeout(header http.Header, defaultTimeout time.Duration) time.Duration {
	if timeout := header.Get(grpcTimeoutHeader); len(timeout) > 0 {
		if t, err := time.ParseDuration(timeout); err == nil {
			return t
		}
	}

	return defaultTimeout
}
