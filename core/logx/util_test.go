package logx

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCaller(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	assert.Contains(t, getCaller(1), filepath.Base(file))
}

func TestGetTimestamp(t *testing.T) {
	ts := getTimestamp()
	tm, err := time.Parse(timeFormat, ts)
	assert.Nil(t, err)
	assert.True(t, time.Since(tm) < time.Minute)
}

func BenchmarkGetCaller(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		getCaller(1)
	}
}
