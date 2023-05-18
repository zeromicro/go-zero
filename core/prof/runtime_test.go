package prof

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDisplayStats(t *testing.T) {
	defer func(orig *os.File) {
		os.Stdout = orig
	}(os.Stdout)

	r, w, _ := os.Pipe()
	os.Stdout = w
	DisplayStats(time.Millisecond * 10)
	time.Sleep(time.Millisecond * 50)
	_ = w.Close()
	out, _ := io.ReadAll(r)
	assert.Contains(t, string(out), "Goroutines: ")
}
