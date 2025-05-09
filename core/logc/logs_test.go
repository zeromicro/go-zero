package logc

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestAddGlobalFields(t *testing.T) {
	buf := logtest.NewCollector(t)

	Info(context.Background(), "hello")
	buf.Reset()

	AddGlobalFields(Field("a", "1"), Field("b", "2"))
	AddGlobalFields(Field("c", "3"))
	Info(context.Background(), "world")
	var m map[string]any
	assert.NoError(t, json.Unmarshal(buf.Bytes(), &m))
	assert.Equal(t, "1", m["a"])
	assert.Equal(t, "2", m["b"])
	assert.Equal(t, "3", m["c"])
}

func TestAlert(t *testing.T) {
	buf := logtest.NewCollector(t)
	Alert(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), "foo"), buf.String())
}

func TestError(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Error(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestErrorf(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Errorf(context.Background(), "foo %s", "bar")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestErrorfn(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Errorfn(context.Background(), func() any {
		return fmt.Sprintf("foo %s", "bar")
	})
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestErrorv(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Errorv(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestErrorw(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Errorw(context.Background(), "foo", Field("a", "b"))
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestInfo(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Info(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestInfof(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Infof(context.Background(), "foo %s", "bar")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestInfofn(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Infofn(context.Background(), func() any {
		return fmt.Sprintf("foo %s", "bar")
	})
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestInfov(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Infov(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestInfow(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Infow(context.Background(), "foo", Field("a", "b"))
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestDebug(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Debug(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestDebugf(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Debugf(context.Background(), "foo %s", "bar")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestDebugfn(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Debugfn(context.Background(), func() any {
		return fmt.Sprintf("foo %s", "bar")
	})
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestDebugv(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Debugv(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestDebugw(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Debugw(context.Background(), "foo", Field("a", "b"))
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)))
}

func TestMust(t *testing.T) {
	assert.NotPanics(t, func() {
		Must(nil)
	})
	assert.NotPanics(t, func() {
		MustSetup(LogConf{})
	})
}

func TestMisc(t *testing.T) {
	SetLevel(logx.DebugLevel)
	assert.NoError(t, SetUp(LogConf{}))
	assert.NoError(t, Close())
}

func TestSlow(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Slow(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)), buf.String())
}

func TestSlowf(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Slowf(context.Background(), "foo %s", "bar")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)), buf.String())
}

func TestSlowfn(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Slowfn(context.Background(), func() any {
		return fmt.Sprintf("foo %s", "bar")
	})
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)), buf.String())
}

func TestSlowv(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Slowv(context.Background(), "foo")
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)), buf.String())
}

func TestSloww(t *testing.T) {
	buf := logtest.NewCollector(t)
	file, line := getFileLine()
	Sloww(context.Background(), "foo", Field("a", "b"))
	assert.True(t, strings.Contains(buf.String(), fmt.Sprintf("%s:%d", file, line+1)), buf.String())
}

func getFileLine() (string, int) {
	_, file, line, _ := runtime.Caller(1)
	short := file

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	return short, line
}
