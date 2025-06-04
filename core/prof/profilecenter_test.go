package prof

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	assert.NotContains(t, generateReport(), "foo")
	report("foo", time.Second)
	assert.Contains(t, generateReport(), "foo")
	report("foo", time.Second)
}
