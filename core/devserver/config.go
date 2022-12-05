package devserver

const (
	defaultPort       = 6470
	defaultMetricPath = "/metrics"
	defaultHealthPath = "/healthz"
)

// Config is config for inner http server.
type Config struct {
	Host         string `json:",optional"`
	Port         int    `json:",default=6470"`
	MetricPath   string `json:",default=/metrics"`
	HealthPath   string `json:",default=/healthz"`
	EnableMetric bool   `json:",default=true"`
	EnablePprof  bool   `json:",optional"`
}

func (c *Config) fillDefault() {
	if *c == (Config{}) {
		c.Port = defaultPort
		c.EnableMetric = true
		c.MetricPath = defaultMetricPath
		c.HealthPath = defaultHealthPath
	}
}
