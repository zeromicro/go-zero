package internal

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const content = "foo"

func TestLoggerError(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Error(content)
	assert.Contains(t, builder.String(), content)
}

func TestLoggerErrorf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Errorf(content)
	assert.Contains(t, builder.String(), content)
}

func TestLoggerErrorln(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Errorln(content)
	assert.Contains(t, builder.String(), content)
}

func TestLoggerFatal(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Fatal(content)
	assert.Contains(t, builder.String(), content)
}

func TestLoggerFatalf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Fatalf(content)
	assert.Contains(t, builder.String(), content)
}

func TestLoggerFatalln(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Fatalln(content)
	assert.Contains(t, builder.String(), content)
}

func TestLoggerWarning(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Warning(content)
	assert.Empty(t, builder.String())
}

func TestLoggerWarningf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Warningf(content)
	assert.Empty(t, builder.String())
}

func TestLoggerWarningln(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	logger := new(Logger)
	logger.Warningln(content)
	assert.Empty(t, builder.String())
}
