package trace

type spanContext struct {
	traceId string
	spanId  string
}

func (sc spanContext) TraceId() string {
	return sc.traceId
}

func (sc spanContext) SpanId() string {
	return sc.spanId
}

func (sc spanContext) Visit(fn func(key, val string) bool) {
	fn(traceIdKey, sc.traceId)
	fn(spanIdKey, sc.spanId)
}
