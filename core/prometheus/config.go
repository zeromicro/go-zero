package prometheus

// A Config is a prometheus config.
type Config struct {
	Host string `json:"Host,optional" yaml:"Host"`
	Port int    `json:"Port,default=9101" yaml:"Port"`
	Path string `json:"Path,default=/metrics" yaml:"Path"`
}
