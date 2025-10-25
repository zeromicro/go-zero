package internal

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	metadataHeaderPrefix = "Grpc-Metadata-"
	metadataPrefix       = "gateway-"
)

// OpenTelemetry trace propagation headers that need to be forwarded to gRPC metadata.
// These headers are used by the W3C Trace Context standard for distributed tracing.
var traceHeaders = map[string]bool{
	"traceparent": true,
	"tracestate":  true,
	"baggage":     true,
}

// ProcessHeaders builds the headers for the gateway from HTTP headers.
// It forwards both custom metadata headers (with Grpc-Metadata- prefix)
// and OpenTelemetry trace propagation headers (traceparent, tracestate, baggage)
// to ensure distributed tracing works correctly across the gateway.
func ProcessHeaders(header http.Header) []string {
	var headers []string

	for k, v := range header {
		// Forward OpenTelemetry trace propagation headers
		// These must be lowercase per gRPC metadata conventions
		if lowerKey := strings.ToLower(k); traceHeaders[lowerKey] {
			for _, vv := range v {
				headers = append(headers, lowerKey+":"+vv)
			}
			continue
		}

		// Forward custom metadata headers with Grpc-Metadata- prefix
		if !strings.HasPrefix(k, metadataHeaderPrefix) {
			continue
		}

		// gRPC metadata keys are case-insensitive and stored as lowercase,
		// so we lowercase the key to match gRPC conventions
		trimmedKey := strings.TrimPrefix(k, metadataHeaderPrefix)
		key := strings.ToLower(fmt.Sprintf("%s%s", metadataPrefix, trimmedKey))
		for _, vv := range v {
			headers = append(headers, key+":"+vv)
		}
	}

	return headers
}
