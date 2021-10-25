package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string `json:",optional"`
	Port int    `json:",default=9101"`
	Path string `json:",default=/metrics"`
}

func (c *Config) Available() bool {
	if c.Port == 0 {
		return false
	}
	return true
}
