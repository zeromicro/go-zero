package trace

// TraceName represents the tracing name.
const TraceName = "go-zero"

// A Config is a open-telemetry config.
type Config struct {
	Name     string  `json:"name,optional" yaml:"Name"`
	Endpoint string  `json:"endpoint,optional" yaml:"Endpoint"`
	Sampler  float64 `json:"sampler,default=1.0" yaml:"Sampler"`
	Batcher  string  `json:"batcher,default=jaeger,options=jaeger|zipkin|grpc" yaml:"Batcher"`
}
