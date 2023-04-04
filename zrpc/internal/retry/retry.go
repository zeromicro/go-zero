package retry

import (
	"fmt"
	"strings"
)

var (
	defaultRetry = `{
	  "name": [{"service": ""}],
	  "retryPolicy": {
		  "MaxAttempts": 3,
		  "InitialBackoff": ".1s",
		  "MaxBackoff": ".1s",
		  "BackoffMultiplier": 1.5,
		  "RetryableStatusCodes": ["UNAVAILABLE", "ABORTED"]
	  }
	}`
)

// MergeRetryPolicy merge user config and defaultRetry config
func MergeRetryPolicy(config string) string {
	c := strings.TrimSuffix(strings.TrimPrefix(config, "["), "]")
	s := ""
	if c != "" {
		s = ","
	}

	return fmt.Sprintf("[%s%s%s]", c, s, defaultRetry)
}
