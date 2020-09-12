package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/tal-tech/go-zero/core/proc"
)

type (
	CounterVecOpts VectorOpts

	CounterVec interface {
		Inc(lables ...string)
		Add(v float64, labels ...string)
		close() bool
	}

	promCounterVec struct {
		counter *prom.CounterVec
	}
)

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

func (cv *promCounterVec) Inc(labels ...string) {
	cv.counter.WithLabelValues(labels...).Inc()
}

func (cv *promCounterVec) Add(v float64, lables ...string) {
	cv.counter.WithLabelValues(lables...).Add(v)
}

func (cv *promCounterVec) close() bool {
	return prom.Unregister(cv.counter)
}
