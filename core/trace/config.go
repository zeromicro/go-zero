package trace

import (
	"strings"
)

// TraceName represents the tracing name.
const TraceName = "go-zero"

// A Config is an opentelemetry config.
type Config struct {
	Name     string  `json:",optional"`
	Endpoint string  `json:",optional"`
	Sampler  float64 `json:",default=1.0"`
	Batcher  string  `json:",default=jaeger,options=jaeger|jaegerudp|zipkin|otlpgrpc|otlphttp"`
}

func (c *Config) getEndpointHost() string {
	EndpointSlice := strings.Split(c.Endpoint, ":")
	if len(EndpointSlice) > 0 {
		return strings.TrimSpace(EndpointSlice[0])
	}
	return ""
}

func (c *Config) getEndpointPort() string {
	EndpointSlice := strings.Split(c.Endpoint, ":")
	if len(EndpointSlice) > 1 {
		return strings.TrimSpace(EndpointSlice[1])
	}

	return ""
}
