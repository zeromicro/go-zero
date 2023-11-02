package metric

import "github.com/zeromicro/go-zero/core/prometheus"

// A VectorOpts is a general configuration.
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

func update(fn func()) {
	if !prometheus.Enabled() {
		return
	}

	fn()
}
