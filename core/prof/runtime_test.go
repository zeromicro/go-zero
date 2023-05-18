package prof

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/internal/iotest"
)

func TestDisplayStats(t *testing.T) {
	iotest.RunTest(t, func() {
		DisplayStats(time.Millisecond * 10)
		time.Sleep(time.Millisecond * 50)
	}, func(s string) {
		assert.Contains(t, s, "Goroutines: ")
	})
}
