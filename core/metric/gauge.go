package metric

import (
	"zero/core/proc"

	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	GaugeVecOpts VectorOpts

	GuageVec interface {
		Set(v float64, labels ...string)
		Inc(labels ...string)
		Add(v float64, labels ...string)
		close() bool
	}

	promGuageVec struct {
		gauge *prom.GaugeVec
	}
)

func NewGaugeVec(cfg *GaugeVecOpts) GuageVec {
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
