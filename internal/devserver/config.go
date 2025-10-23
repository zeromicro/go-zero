package devserver

// Config is config for inner http server.
type Config struct {
	Enabled        bool   `json:",default=true"`
	Host           string `json:",optional"`
	Port           int    `json:",default=6060"`
	MetricsPath    string `json:",default=/metrics"`
	HealthPath     string `json:",default=/healthz"`
	EnableMetrics  bool   `json:",default=true"`
	EnablePprof    bool   `json:",default=true"`
	HealthResponse string `json:",default=OK"`
	// EnableOpenMetrics expose promhttp.HandlerOpts.EnableOpenMetrics option to add the experimental OpenMetrics encoding,
	// which is the only way to transmit exemplars.
	EnableOpenMetrics bool `json:",optional"`
}
