package trace

// TraceName represents the tracing name.
const TraceName = "go-zero"

// A Config is a open-telemetry config.
type Config struct {
	Name     string  `json:"Name,optional" yaml:"Name"`
	Endpoint string  `json:"Endpoint,optional" yaml:"Endpoint"`
	Sampler  float64 `json:"Sampler,default=1.0" yaml:"Sampler"`
	Batcher  string  `json:"Batcher,default=jaeger,options=jaeger|zipkin|grpc" yaml:"Batcher"`
}
