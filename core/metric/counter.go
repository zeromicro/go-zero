package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/proc"
)

type (
	// A CounterVecOpts is an alias of VectorOpts.
	CounterVecOpts VectorOpts

	// CounterVec interface represents a counter vector.
	CounterVec interface {
		// Inc increments labels.
		Inc(labels ...string)
		// Add adds labels with v.
		Add(v float64, labels ...string)
		close() bool
	}

	promCounterVec struct {
		counter *prom.CounterVec
	}
)

// NewCounterVec returns a CounterVec.
func NewCounterVec(cfg *CounterVecOpts) CounterVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewCounterVec(prom.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels)
	prom.MustRegister(vec)
	cv := &promCounterVec{
		counter: vec,
	}
	proc.AddShutdownListener(func() {
		cv.close()
	})

	return cv
}

func (cv *promCounterVec) Add(v float64, labels ...string) {
	update(func() {
		cv.counter.WithLabelValues(labels...).Add(v)
	})
}

func (cv *promCounterVec) Inc(labels ...string) {
	update(func() {
		cv.counter.WithLabelValues(labels...).Inc()
	})
}

func (cv *promCounterVec) close() bool {
	return prom.Unregister(cv.counter)
}
