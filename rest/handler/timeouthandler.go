package handler

import (
	"net/http"
	"time"
)

const reason = "Request Timeout"

func TimeoutHandler(duration time.Duration, msg string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration > 0 {
			if msg == "" {
				msg = reason
			}
			return http.TimeoutHandler(next, duration, msg)
		} else {
			return next
		}
	}
}
