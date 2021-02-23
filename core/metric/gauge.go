package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/tal-tech/go-zero/core/proc"
)

type (
	// GaugeVecOpts is an alias of VectorOpts.
	GaugeVecOpts VectorOpts

	// GaugeVec represents a gauge vector.
	GaugeVec interface {
		// Set sets v to labels.
		Set(v float64, labels ...string)
		// Inc increments labels.
		Inc(labels ...string)
		// Add adds v to labels.
		Add(v float64, labels ...string)
		close() bool
	}

	promGuageVec struct {
		gauge *prom.GaugeVec
	}
)

// NewGaugeVec returns a GaugeVec.
func NewGaugeVec(cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)
	prom.MustRegister(vec)
	gv := &promGuageVec{
		gauge: vec,
	}
	proc.AddShutdownListener(func() {
		gv.close()
	})

	return gv
}

func (gv *promGuageVec) Inc(labels ...string) {
	gv.gauge.WithLabelValues(labels...).Inc()
}

func (gv *promGuageVec) Add(v float64, lables ...string) {
	gv.gauge.WithLabelValues(lables...).Add(v)
}

func (gv *promGuageVec) Set(v float64, lables ...string) {
	gv.gauge.WithLabelValues(lables...).Set(v)
}

func (gv *promGuageVec) close() bool {
	return prom.Unregister(gv.gauge)
}
