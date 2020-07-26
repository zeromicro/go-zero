package metric

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewHistogramVec(t *testing.T) {
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Name:    "duration_ms",
		Help:    "rpc server requests duration(ms).",
		Buckets: []float64{1, 2, 3},
	})
	defer histogramVec.close()
	histogramVecNil := NewHistogramVec(nil)
	assert.NotNil(t, histogramVec)
	assert.Nil(t, histogramVecNil)
}

func TestHistogramObserve(t *testing.T) {
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Name:    "counts",
		Help:    "rpc server requests duration(ms).",
		Buckets: []float64{1, 2, 3},
		Labels:  []string{"method"},
	})
	defer histogramVec.close()
	hv, _ := histogramVec.(*promHistogramVec)
	hv.Observe(2, "/Users")

	metadata := `
		# HELP counts rpc server requests duration(ms).
        # TYPE counts histogram
`
	val := `
		counts_bucket{method="/Users",le="1"} 0
		counts_bucket{method="/Users",le="2"} 1
		counts_bucket{method="/Users",le="3"} 1
		counts_bucket{method="/Users",le="+Inf"} 1
		counts_sum{method="/Users"} 2
        counts_count{method="/Users"} 1
`

	err := testutil.CollectAndCompare(hv.histogram, strings.NewReader(metadata+val))
	assert.Nil(t, err)
}
