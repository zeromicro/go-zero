package prof

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDisplayStats(t *testing.T) {
	var buf strings.Builder
	displayStatsWithWriter(&buf, time.Millisecond*10)
	time.Sleep(time.Millisecond * 50)
	assert.Contains(t, buf.String(), "Goroutines: ")
}
