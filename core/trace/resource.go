package trace

import "go.opentelemetry.io/otel/attribute"

var attrResources = make([]attribute.KeyValue, 0)

func AddResources(attrs ...attribute.KeyValue) {
	attrResources = append(attrResources, attrs...)
}
