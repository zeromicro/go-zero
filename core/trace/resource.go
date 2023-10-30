package trace

import "go.opentelemetry.io/otel/attribute"

var attrResources = make([]attribute.KeyValue, 0)

// AddResources add more resources in addition to configured trace name.
func AddResources(attrs ...attribute.KeyValue) {
	attrResources = append(attrResources, attrs...)
}
