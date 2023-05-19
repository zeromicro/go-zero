package prof

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDisplayStats(t *testing.T) {
	writer := &threadSafeBuffer{
		buf: strings.Builder{},
	}
	displayStatsWithWriter(writer, time.Millisecond*10)
	time.Sleep(time.Millisecond * 50)
	assert.Contains(t, writer.String(), "Goroutines: ")
}

type threadSafeBuffer struct {
	buf  strings.Builder
	lock sync.Mutex
}

func (b *threadSafeBuffer) String() string {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.buf.String()
}

func (b *threadSafeBuffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.buf.Write(p)
}
