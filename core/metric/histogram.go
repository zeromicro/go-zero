package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/tal-tech/go-zero/core/proc"
)

type (
	HistogramVecOpts struct {
		Namespace string
		Subsystem string
		Name      string
		Help      string
		Labels    []string
		Buckets   []float64
	}

	HistogramVec interface {
		Observe(v int64, lables ...string)
		close() bool
	}

	promHistogramVec struct {
		histogram *prom.HistogramVec
	}
)

func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewHistogramVec(prom.HistogramOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
		Buckets:   cfg.Buckets,
	}, cfg.Labels)
	prom.MustRegister(vec)
	hv := &promHistogramVec{
		histogram: vec,
	}
	proc.AddShutdownListener(func() {
		hv.close()
	})

	return hv
}

func (hv *promHistogramVec) Observe(v int64, labels ...string) {
	hv.histogram.WithLabelValues(labels...).Observe(float64(v))
}

func (hv *promHistogramVec) close() bool {
	return prom.Unregister(hv.histogram)
}
