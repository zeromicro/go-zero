package trace

import "net/http"

// TraceIdKey is the trace id header.
// https://www.w3.org/TR/trace-context/#trace-id
// May change it to trace-id afterwards.
var TraceIdKey = http.CanonicalHeaderKey("x-trace-id")
