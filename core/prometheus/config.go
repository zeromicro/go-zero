package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string `json:"host,optional" yaml:"Host"`
	Port int    `json:"port,default=9101" yaml:"Port"`
	Path string `json:"path,default=/metrics" yaml:"Path"`
}
