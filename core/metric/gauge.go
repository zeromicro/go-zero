package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/proc"
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
		// Dec decrements labels.
		Dec(labels ...string)
		// Add adds v to labels.
		Add(v float64, labels ...string)
		// Sub subtracts v to labels.
		Sub(v float64, labels ...string)
		close() bool
	}

	promGaugeVec struct {
		gauge *prom.GaugeVec
	}
)

// NewGaugeVec returns a GaugeVec.
func NewGaugeVec(cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewGaugeVec(prom.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels)
	prom.MustRegister(vec)
	gv := &promGaugeVec{
		gauge: vec,
	}
	proc.AddShutdownListener(func() {
		gv.close()
	})

	return gv
}

func (gv *promGaugeVec) Add(v float64, labels ...string) {
	update(func() {
		gv.gauge.WithLabelValues(labels...).Add(v)
	})
}

func (gv *promGaugeVec) Dec(labels ...string) {
	update(func() {
		gv.gauge.WithLabelValues(labels...).Dec()
	})
}

func (gv *promGaugeVec) Inc(labels ...string) {
	update(func() {
		gv.gauge.WithLabelValues(labels...).Inc()
	})
}

func (gv *promGaugeVec) Set(v float64, labels ...string) {
	update(func() {
		gv.gauge.WithLabelValues(labels...).Set(v)
	})
}

func (gv *promGaugeVec) Sub(v float64, labels ...string) {
	update(func() {
		gv.gauge.WithLabelValues(labels...).Sub(v)
	})
}

func (gv *promGaugeVec) close() bool {
	return prom.Unregister(gv.gauge)
}
