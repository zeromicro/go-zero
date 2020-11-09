package prometheus

type Config struct {
	Host string `json:",default=127.0.0.1"`
	Port int    `json:",default=9101"`
	Path string `json:",default=/metrics"`
}
