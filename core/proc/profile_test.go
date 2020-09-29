package proc

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	var buf strings.Builder
	log.SetOutput(&buf)
	profiler := StartProfile()
	// start again should not work
	assert.NotNil(t, StartProfile())
	profiler.Stop()
	// stop twice
	profiler.Stop()
	assert.True(t, strings.Contains(buf.String(), ".pprof"))
}
