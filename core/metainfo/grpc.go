package metainfo

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

var _ propagation.TextMapCarrier = (*GrpcHeaderCarrier)(nil)

// GrpcHeaderCarrier implements propagation.TextMapCarrier for grpc metadata.MD.
type GrpcHeaderCarrier metadata.MD

// Get returns the value associated with the passed key.
func (mc GrpcHeaderCarrier) Get(key string) string {
	vals := mc[key]
	if len(vals) == 0 {
		return ""
	}

	return vals[0]
}

// Keys lists the keys stored in this carrier.
func (mc GrpcHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}

// Set stores the key-value pair.
func (mc GrpcHeaderCarrier) Set(key, value string) {
	metadata.MD(mc).Set(key, value)
}
