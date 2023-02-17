package trace

import "fmt"

// TraceName represents the tracing name.
const TraceName = "go-zero"

// A Config is an opentelemetry config.
type Config struct {
	Name      string  `json:",optional"`
	AgentHost string  `json:",optional"`
	AgentPort string  `json:",optional"`
	Endpoint  string  `json:",optional"`
	Sampler   float64 `json:",default=1.0"`
	Batcher   string  `json:",default=jaeger,options=jaeger|jaegerudp|zipkin|otlpgrpc|otlphttp"`
}

func (c *Config) isAgentEndPoint() bool {
	return len(c.AgentHost) != 0 && len(c.AgentPort) != 0
}

func (c *Config) getEndpoint() string {
	if c.isAgentEndPoint() {
		return fmt.Sprintf("%s:%s", c.AgentHost, c.AgentPort)
	}
	return c.Endpoint
}
