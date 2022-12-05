package health

import (
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/syncx"
)

// defaultHealthManager is global comboHealthManager for byone self.
var defaultHealthManager = newComboHealthManager()

type (
	// Probe represents readiness status of given component.
	Probe interface {
		// MarkReady sets a ready state for the endpoint handlers.
		MarkReady()
		// MarkNotReady sets a not ready state for the endpoint handlers.
		MarkNotReady()
		// IsReady return inner state for the component.
		IsReady() bool
		// Name return probe name identifier
		Name() string
	}

	// healthManager manage app healthy.
	healthManager struct {
		ready syncx.AtomicBool
		name  string
	}

	// comboHealthManager folds given probes into one, reflects their statuses in a thread-safe way.
	comboHealthManager struct {
		mu     sync.Mutex
		probes []Probe
	}
)

// AddProbe add components probe to global comboHealthManager.
func AddProbe(probe Probe) {
	defaultHealthManager.addProbe(probe)
}

// NewHealthManager returns a new healthManager.
func NewHealthManager(name string) Probe {
	return &healthManager{
		name: name,
	}
}

// MarkReady sets a ready state for the endpoint handlers.
func (h *healthManager) MarkReady() {
	h.ready.Set(true)
}

// MarkNotReady sets a not ready state for the endpoint handlers.
func (h *healthManager) MarkNotReady() {
	h.ready.Set(false)
}

// IsReady return inner state for the component.
func (h *healthManager) IsReady() bool {
	return h.ready.True()
}

// Name return probe name identifier
func (h *healthManager) Name() string {
	return h.name
}

func newComboHealthManager() *comboHealthManager {
	return &comboHealthManager{}
}

// MarkReady sets components status to ready.
func (p *comboHealthManager) MarkReady() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, probe := range p.probes {
		probe.MarkReady()
	}
}

// MarkNotReady sets components status to not ready with given error as a cause.
func (p *comboHealthManager) MarkNotReady() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, probe := range p.probes {
		probe.MarkNotReady()
	}
}

// IsReady return composed status of all components.
func (p *comboHealthManager) IsReady() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, probe := range p.probes {
		if !probe.IsReady() {
			return false
		}
	}
	return true
}

func (p *comboHealthManager) verboseInfo() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	var info string
	for _, probe := range p.probes {
		if probe.IsReady() {
			info += probe.Name() + " is ready; "
		} else {
			info += probe.Name() + " is not ready; "
		}
	}
	return info
}

// addProbe add components probe to comboHealthManager.
func (p *comboHealthManager) addProbe(probe Probe) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.probes = append(p.probes, probe)
}

// CreateHttpHandler create health http handler base on given probe.
func CreateHttpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		if defaultHealthManager.IsReady() {
			_, _ = w.Write([]byte("OK"))
		} else {
			http.Error(w, "Service Unavailable\n"+defaultHealthManager.verboseInfo(), http.StatusServiceUnavailable)
		}
	}
}
