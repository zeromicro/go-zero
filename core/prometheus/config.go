package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string `json:",optional"`
	Port int    `json:",default=9101"`
	Path string `json:",default=/metrics"`
}

type MetricsPusherConfig struct {
	Url      string
	JobName  string
	Interval int `json:",default=300"`
}
