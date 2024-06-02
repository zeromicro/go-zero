package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/zeromicro/go-zero/core/proc"
)

type (
	// A SummaryVecOpts is a summary vector options
	SummaryVecOpts struct {
		VecOpt     VectorOpts
		Objectives map[float64]float64
	}

	// A SummaryVec interface represents a summary vector.
	SummaryVec interface {
		// Observe adds observation v to labels.
		Observe(v float64, labels ...string)
		close() bool
	}

	promSummaryVec struct {
		summary *prom.SummaryVec
	}
)

// NewSummaryVec return a SummaryVec
func NewSummaryVec(cfg *SummaryVecOpts) SummaryVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewSummaryVec(
		prom.SummaryOpts{
			Namespace:  cfg.VecOpt.Namespace,
			Subsystem:  cfg.VecOpt.Subsystem,
			Name:       cfg.VecOpt.Name,
			Help:       cfg.VecOpt.Help,
			Objectives: cfg.Objectives,
		},
		cfg.VecOpt.Labels,
	)
	prom.MustRegister(vec)
	sv := &promSummaryVec{
		summary: vec,
	}
	proc.AddShutdownListener(func() {
		sv.close()
	})

	return sv
}

func (sv *promSummaryVec) Observe(v float64, labels ...string) {
	update(func() {
		sv.summary.WithLabelValues(labels...).Observe(v)
	})
}

func (sv *promSummaryVec) close() bool {
	return prom.Unregister(sv.summary)
}
