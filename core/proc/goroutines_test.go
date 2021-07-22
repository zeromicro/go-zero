package proc

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDumpGoroutines(t *testing.T) {
	var buf strings.Builder
	log.SetOutput(&buf)
	dumpGoroutines()
	assert.True(t, strings.Contains(buf.String(), ".dump"))
}
