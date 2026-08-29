package metric

import (
	"regexp"
	"strings"
	"testing"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

func TestNewHistogramVec(t *testing.T) {
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Name:    "duration_ms",
		Help:    "rpc server requests duration(ms).",
		Buckets: []float64{1, 2, 3},
	})
	defer histogramVec.(*promHistogramVec).close()
	histogramVecNil := NewHistogramVec(nil)
	assert.NotNil(t, histogramVec)
	assert.Nil(t, histogramVecNil)
}

func TestHistogramObserve(t *testing.T) {
	startAgent()
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Name:    "counts",
		Help:    "rpc server requests duration(ms).",
		Buckets: []float64{1, 2, 3},
		Labels:  []string{"method"},
	})
	defer histogramVec.(*promHistogramVec).close()
	hv, _ := histogramVec.(*promHistogramVec)
	hv.Observe(2, "/Users")
	hv.ObserveFloat(1.1, "/Users")

	metadata := `
		# HELP counts rpc server requests duration(ms).
        # TYPE counts histogram
`
	val := `
		counts_bucket{method="/Users",le="1"} 0
		counts_bucket{method="/Users",le="2"} 2
		counts_bucket{method="/Users",le="3"} 2
		counts_bucket{method="/Users",le="+Inf"} 2
		counts_sum{method="/Users"} 3.1
        counts_count{method="/Users"} 2
`

	err := testutil.CollectAndCompare(hv.histogram, strings.NewReader(metadata+val))
	assert.Nil(t, err)
}

func Test_promHistogramVec_ObserveWithExemplar(t *testing.T) {
	startAgent()
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Name:    "counts",
		Help:    "rpc server requests duration(ms).",
		Buckets: []float64{1, 2, 3},
		Labels:  []string{"method"},
	})
	defer histogramVec.(*promHistogramVec).close()

	histogramVec.ObserveWithExemplar(1.5, prom.Labels{"test": "test15"}, "/Users")
	histogramVec.ObserveWithExemplar(2.5, prom.Labels{"test": "test25"}, "/Users")
	histogramVec.ObserveWithExemplar(3.5, prom.Labels{"test": "test35"}, "/Users")
	hv, _ := histogramVec.(*promHistogramVec)

	expect := `# HELP counts rpc server requests duration(ms).
# TYPE counts histogram
counts_bucket{method="/Users",le="1.0"} 0
counts_bucket{method="/Users",le="2.0"} 1 # {test="test15"} 1.5
counts_bucket{method="/Users",le="3.0"} 2 # {test="test25"} 2.5
counts_bucket{method="/Users",le="+Inf"} 3 # {test="test35"} 3.5
counts_sum{method="/Users"} 7.5
counts_count{method="/Users"} 3
`
	m, err := testutil.CollectAndFormat(hv.histogram, expfmt.TypeOpenMetrics, "counts")
	assert.NoError(t, err)
	assert.Equal(t, expect, removeTimestamp(string(m)))
}

func Test_promHistogramVec_ObserveWithTrace(t *testing.T) {
	startAgent()
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Name:    "counts",
		Help:    "rpc server requests duration(ms).",
		Buckets: []float64{1, 2, 3},
		Labels:  []string{"method"},
	})
	defer histogramVec.(*promHistogramVec).close()

	histogramVec.ObserveWithTrace(1.5, "4bf92f3577b34da6a3ce929d0e4e6b3a", "/Users")
	histogramVec.ObserveWithTrace(2.5, "1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d", "/Users")
	histogramVec.ObserveWithTrace(3.5, "8e7f6a5d4c3b2a1f0e9d8c7b6a5d4c3b", "/Users")
	hv, _ := histogramVec.(*promHistogramVec)

	expect := `# HELP counts rpc server requests duration(ms).
# TYPE counts histogram
counts_bucket{method="/Users",le="1.0"} 0
counts_bucket{method="/Users",le="2.0"} 1 # {trace_id="4bf92f3577b34da6a3ce929d0e4e6b3a"} 1.5
counts_bucket{method="/Users",le="3.0"} 2 # {trace_id="1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d"} 2.5
counts_bucket{method="/Users",le="+Inf"} 3 # {trace_id="8e7f6a5d4c3b2a1f0e9d8c7b6a5d4c3b"} 3.5
counts_sum{method="/Users"} 7.5
counts_count{method="/Users"} 3
`
	m, err := testutil.CollectAndFormat(hv.histogram, expfmt.TypeOpenMetrics, "counts")
	assert.NoError(t, err)
	assert.Equal(t, expect, removeTimestamp(string(m)))
}

// removeTimestamp removes the timestamp from the OpenMetrics output,
// eg: counts_bucket{method="/Users",le="2.0"} 1 # {test="test15"} 1.5 1.7442025686415942e+09
func removeTimestamp(s string) string {
	r := regexp.MustCompile(`\s+\d+\.\d+e[+-]\d+`)
	return r.ReplaceAllString(s, "")
}
