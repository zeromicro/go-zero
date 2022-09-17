package trace

// TraceName represents the tracing name.
const TraceName = "go-zero"

// A Config is a opentelemetry config.
type Config struct {
	Name     string  `json:",optional"`
	Endpoint string  `json:",optional"`
	Sampler  float64 `default:"1.0"`
	Batcher  string  `default:"jaeger" validate:"oneof=jaeger zipkin grpc"`
}
