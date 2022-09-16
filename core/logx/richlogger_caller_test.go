package logx

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogWithCallerSkip(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	l := WithCallerSkip(1).WithCallerSkip(0)
	p := func(v string) {
		l.Info(v)
	}

	p(testlog)
	assert.True(t, strings.Contains(w.String(), "logx/richlogger_caller_test.go:25"), w.String())

	w.Reset()
	l = WithCallerSkip(0).WithCallerSkip(1)
	p(testlog)
	assert.True(t, strings.Contains(w.String(), "logx/richlogger_caller_test.go:30"), w.String())
}

func TestWithFields(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	l := WithFields(Field("a", "b")).WithFields(Field("c", "d"))
	p := func(v string) {
		l.Info(v)
	}

	p(testlog)
	var v struct {
		A string `json:"a"`
		C string `json:"c"`
	}
	assert.Nil(t, json.Unmarshal([]byte(w.String()), &v))
	assert.Equal(t, "b", v.A, w.String())
	assert.Equal(t, "d", v.C, w.String())
}
