package prof

import (
	"testing"

	"github.com/l306287405/go-zero/core/utils"
)

func TestProfiler(t *testing.T) {
	EnableProfiling()
	Start()
	Report("foo", ProfilePoint{
		ElapsedTimer: utils.NewElapsedTimer(),
	})
}

func TestNullProfiler(t *testing.T) {
	p := newNullProfiler()
	p.Start()
	p.Report("foo", ProfilePoint{
		ElapsedTimer: utils.NewElapsedTimer(),
	})
}
