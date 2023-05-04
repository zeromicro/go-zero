package proc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestDumpGoroutines(t *testing.T) {
	buf := logtest.NewCollector(t)
	dumpGoroutines()
	assert.True(t, strings.Contains(buf.String(), ".dump"))
}
