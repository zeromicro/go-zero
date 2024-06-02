package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

const content = "foo"

func TestLoggerError(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Error(content)
	assert.Contains(t, c.String(), content)
}

func TestLoggerErrorf(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Errorf(content)
	assert.Contains(t, c.String(), content)
}

func TestLoggerErrorln(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Errorln(content)
	assert.Contains(t, c.String(), content)
}

func TestLoggerFatal(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Fatal(content)
	assert.Contains(t, c.String(), content)
}

func TestLoggerFatalf(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Fatalf(content)
	assert.Contains(t, c.String(), content)
}

func TestLoggerFatalln(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Fatalln(content)
	assert.Contains(t, c.String(), content)
}

func TestLoggerInfo(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Info(content)
	assert.Empty(t, c.String())
}

func TestLoggerInfof(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Infof(content)
	assert.Empty(t, c.String())
}

func TestLoggerWarning(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Warning(content)
	assert.Empty(t, c.String())
}

func TestLoggerInfoln(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Infoln(content)
	assert.Empty(t, c.String())
}

func TestLoggerWarningf(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Warningf(content)
	assert.Empty(t, c.String())
}

func TestLoggerWarningln(t *testing.T) {
	c := logtest.NewCollector(t)
	logger := new(Logger)
	logger.Warningln(content)
	assert.Empty(t, c.String())
}

func TestLogger_V(t *testing.T) {
	logger := new(Logger)
	// grpclog.fatalLog
	assert.True(t, logger.V(3))
	// grpclog.infoLog
	assert.False(t, logger.V(0))
}
