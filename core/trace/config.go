package trace

// TraceName represents the tracing name.
const TraceName = "go-zero"

// A Config is an opentelemetry config.
type Config struct {
	Name     string  `json:",optional"`
	Endpoint string  `json:",optional"`
	Sampler  float64 `json:",default=1.0"`
	Batcher  string  `json:",default=jaeger,options=jaeger|zipkin|otlpgrpc|otlphttp|file"`
	// OtlpHeaders represents the headers for OTLP gRPC or HTTP transport.
	// For example:
	//  uptrace-dsn: 'http://project2_secret_token@localhost:14317/2'
	OtlpHeaders map[string]string `json:",optional"`
	// OtlpHttpPath represents the path for OTLP HTTP transport.
	// For example
	// /v1/traces
	OtlpHttpPath string `json:",optional"`
	// Disabled indicates whether StartAgent starts the agent.
	Disabled bool `json:",optional"`
}
