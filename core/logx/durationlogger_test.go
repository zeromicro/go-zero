package logx

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithDurationError(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Error("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationErrorf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Errorf("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationInfo(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Info("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationInfof(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Infof("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationSlow(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Slow("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationSlowf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).WithDuration(time.Hour).Slowf("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}
