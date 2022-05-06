package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

const content = "foo"

func TestLoggerError(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Error(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerErrorf(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Errorf(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerErrorln(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Errorln(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerFatal(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Fatal(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerFatalf(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Fatalf(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerFatalln(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Fatalln(content)
	assert.Contains(t, w.String(), content)
}

func TestLoggerInfo(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Info(content)
	assert.Empty(t, w.String())
}

func TestLoggerInfof(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Infof(content)
	assert.Empty(t, w.String())
}

func TestLoggerWarning(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Warning(content)
	assert.Empty(t, w.String())
}

func TestLoggerInfoln(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Infoln(content)
	assert.Empty(t, w.String())
}

func TestLoggerWarningf(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Warningf(content)
	assert.Empty(t, w.String())
}

func TestLoggerWarningln(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logger := new(Logger)
	logger.Warningln(content)
	assert.Empty(t, w.String())
}

func TestLogger_V(t *testing.T) {
	logger := new(Logger)
	// grpclog.fatalLog
	assert.True(t, logger.V(3))
	// grpclog.infoLog
	assert.False(t, logger.V(0))
}

func injectLog() (r *strings.Builder, restore func()) {
	var buf strings.Builder
	w := logx.NewWriter(&buf)
	o := logx.Reset()
	logx.SetWriter(w)

	return &buf, func() {
		logx.Reset()
		logx.SetWriter(o)
	}
}
