package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/prometheus"
)

func TestNewCounterVec(t *testing.T) {
	counterVec := NewCounterVec(&CounterVecOpts{
		Namespace: "http_server",
		Subsystem: "requests",
		Name:      "total",
		Help:      "rpc client requests error count.",
	})
	defer counterVec.close()
	counterVecNil := NewCounterVec(nil)
	counterVec.Inc("path", "code")
	counterVec.Add(1, "path", "code")
	proc.Shutdown()
	assert.NotNil(t, counterVec)
	assert.Nil(t, counterVecNil)
}

func TestCounterIncr(t *testing.T) {
	startAgent()
	counterVec := NewCounterVec(&CounterVecOpts{
		Namespace: "http_client",
		Subsystem: "call",
		Name:      "code_total",
		Help:      "http client requests error count.",
		Labels:    []string{"path", "code"},
	})
	defer counterVec.close()
	cv, _ := counterVec.(*promCounterVec)
	cv.Inc("/Users", "500")
	cv.Inc("/Users", "500")
	r := testutil.ToFloat64(cv.counter)
	assert.Equal(t, float64(2), r)
}

func TestCounterAdd(t *testing.T) {
	startAgent()
	counterVec := NewCounterVec(&CounterVecOpts{
		Namespace: "rpc_server",
		Subsystem: "requests",
		Name:      "err_total",
		Help:      "rpc client requests error count.",
		Labels:    []string{"method", "code"},
	})
	defer counterVec.close()
	cv, _ := counterVec.(*promCounterVec)
	cv.Add(11, "/Users", "500")
	cv.Add(22, "/Users", "500")
	r := testutil.ToFloat64(cv.counter)
	assert.Equal(t, float64(33), r)
}

func startAgent() {
	prometheus.StartAgent(prometheus.Config{
		Host: "127.0.0.1",
		Port: 9101,
		Path: "/metrics",
	})
}
