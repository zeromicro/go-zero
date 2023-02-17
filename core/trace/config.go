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

func (c *Config) parseEndpoint() (host string, port string) {
	EndpointSlice := strings.Split(c.Endpoint, ":")
	if len(EndpointSlice) > 0 {
		host = strings.TrimSpace(EndpointSlice[0])
	}
	if len(EndpointSlice) > 0 {
		port = strings.TrimSpace(EndpointSlice[1])
	}

	return host, port
}
