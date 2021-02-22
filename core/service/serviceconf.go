package service

import (
	"log"

	"github.com/tal-tech/go-zero/core/load"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/prometheus"
	"github.com/tal-tech/go-zero/core/stat"
)

const (
	// DevMode means development mode.
	DevMode = "dev"
	// TestMode means test mode.
	TestMode = "test"
	// PreMode means pre-release mode.
	PreMode = "pre"
	// ProMode means production mode.
	ProMode = "pro"
)

// A ServiceConf is a service config.
type ServiceConf struct {
	Name       string
	Log        logx.LogConf
	Mode       string            `json:",default=pro,options=dev|test|rt|pre|pro"`
	MetricsUrl string            `json:",optional"`
	Prometheus prometheus.Config `json:",optional"`
}

// MustSetUp sets up the service, exits on error.
func (sc ServiceConf) MustSetUp() {
	if err := sc.SetUp(); err != nil {
		log.Fatal(err)
	}
}

// SetUp sets up the service.
func (sc ServiceConf) SetUp() error {
	if len(sc.Log.ServiceName) == 0 {
		sc.Log.ServiceName = sc.Name
	}
	if err := logx.SetUp(sc.Log); err != nil {
		return err
	}

	sc.initMode()
	prometheus.StartAgent(sc.Prometheus)
	if len(sc.MetricsUrl) > 0 {
		stat.SetReportWriter(stat.NewRemoteWriter(sc.MetricsUrl))
	}

	return nil
}

func (sc ServiceConf) initMode() {
	switch sc.Mode {
	case DevMode, TestMode, PreMode:
		load.Disable()
		stat.SetReporter(nil)
	}
}
