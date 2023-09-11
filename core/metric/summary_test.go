package metric

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
)

func TestNewSummaryVec(t *testing.T) {
	summaryVec := NewSummaryVec(&SummaryVecOpts{
		VecOpt: VectorOpts{
			Namespace: "http_server",
			Subsystem: "requests",
			Name:      "duration_quantiles",
			Help:      "rpc client requests duration(ms) φ quantiles ",
			Labels:    []string{"method"},
		},
		Objectives: map[float64]float64{
			0.5: 0.01,
			0.9: 0.01,
		},
	})
	defer summaryVec.close()
	summaryVecNil := NewSummaryVec(nil)
	assert.NotNil(t, summaryVec)
	assert.Nil(t, summaryVecNil)
}

func TestSummaryObserve(t *testing.T) {
	startAgent()
	summaryVec := NewSummaryVec(&SummaryVecOpts{
		VecOpt: VectorOpts{
			Namespace: "http_server",
			Subsystem: "requests",
			Name:      "duration_quantiles",
			Help:      "rpc client requests duration(ms) φ quantiles ",
			Labels:    []string{"method"},
		},
		Objectives: map[float64]float64{
			0.3: 0.01,
			0.6: 0.01,
			1:   0.01,
		},
	})
	defer summaryVec.close()
	sv := summaryVec.(*promSummaryVec)
	sv.Observe(100, "GET")
	sv.Observe(200, "GET")
	sv.Observe(300, "GET")
	metadata := `
		# HELP http_server_requests_duration_quantiles rpc client requests duration(ms) φ quantiles 
		# TYPE http_server_requests_duration_quantiles summary
`
	val := `
		http_server_requests_duration_quantiles{method="GET",quantile="0.3"} 100
		http_server_requests_duration_quantiles{method="GET",quantile="0.6"} 200
		http_server_requests_duration_quantiles{method="GET",quantile="1"} 300
		http_server_requests_duration_quantiles_sum{method="GET"} 600
		http_server_requests_duration_quantiles_count{method="GET"} 3
`

	err := testutil.CollectAndCompare(sv.summary, strings.NewReader(metadata+val))
	assert.Nil(t, err)
	proc.Shutdown()
}
