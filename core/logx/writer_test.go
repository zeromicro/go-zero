package logx

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewWriter(t *testing.T) {
	const literal = "foo bar"
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Info(literal)
	assert.Contains(t, buf.String(), literal)
	buf.Reset()
	w.Debug(literal)
	assert.Contains(t, buf.String(), literal)
}

func TestConsoleWriter(t *testing.T) {
	var buf bytes.Buffer
	w := newConsoleWriter()
	lw := newLogWriter(log.New(&buf, "", 0))
	w.(*concreteWriter).errorLog = lw
	w.Alert("foo bar 1")
	var val mockedEntry
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelAlert, val.Level)
	assert.Equal(t, "foo bar 1", val.Content)

	buf.Reset()
	w.(*concreteWriter).errorLog = lw
	w.Error("foo bar 2")
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelError, val.Level)
	assert.Equal(t, "foo bar 2", val.Content)

	buf.Reset()
	w.(*concreteWriter).infoLog = lw
	w.Info("foo bar 3")
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelInfo, val.Level)
	assert.Equal(t, "foo bar 3", val.Content)

	buf.Reset()
	w.(*concreteWriter).severeLog = lw
	w.Severe("foo bar 4")
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelFatal, val.Level)
	assert.Equal(t, "foo bar 4", val.Content)

	buf.Reset()
	w.(*concreteWriter).slowLog = lw
	w.Slow("foo bar 5")
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelSlow, val.Level)
	assert.Equal(t, "foo bar 5", val.Content)

	buf.Reset()
	w.(*concreteWriter).statLog = lw
	w.Stat("foo bar 6")
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelStat, val.Level)
	assert.Equal(t, "foo bar 6", val.Content)

	w.(*concreteWriter).infoLog = hardToCloseWriter{}
	assert.NotNil(t, w.Close())
	w.(*concreteWriter).infoLog = easyToCloseWriter{}
	w.(*concreteWriter).errorLog = hardToCloseWriter{}
	assert.NotNil(t, w.Close())
	w.(*concreteWriter).errorLog = easyToCloseWriter{}
	w.(*concreteWriter).severeLog = hardToCloseWriter{}
	assert.NotNil(t, w.Close())
	w.(*concreteWriter).severeLog = easyToCloseWriter{}
	w.(*concreteWriter).slowLog = hardToCloseWriter{}
	assert.NotNil(t, w.Close())
	w.(*concreteWriter).slowLog = easyToCloseWriter{}
	w.(*concreteWriter).statLog = hardToCloseWriter{}
	assert.NotNil(t, w.Close())
	w.(*concreteWriter).statLog = easyToCloseWriter{}
}

func TestNewFileWriter(t *testing.T) {
	t.Run("access", func(t *testing.T) {
		_, err := newFileWriter(LogConf{
			Path: "/not-exists",
		})
		assert.Error(t, err)
	})
}

func TestNopWriter(t *testing.T) {
	assert.NotPanics(t, func() {
		var w nopWriter
		w.Alert("foo")
		w.Debug("foo")
		w.Error("foo")
		w.Info("foo")
		w.Severe("foo")
		w.Stack("foo")
		w.Stat("foo")
		w.Slow("foo")
		_ = w.Close()
	})
}

func TestWriteJson(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	writeJson(nil, "foo")
	assert.Contains(t, buf.String(), "foo")

	buf.Reset()
	writeJson(hardToWriteWriter{}, "foo")
	assert.Contains(t, buf.String(), "write error")

	buf.Reset()
	writeJson(nil, make(chan int))
	assert.Contains(t, buf.String(), "unsupported type")

	buf.Reset()
	type C struct {
		RC func()
	}
	writeJson(nil, C{
		RC: func() {},
	})
	assert.Contains(t, buf.String(), "runtime/debug.Stack")
}

func TestWritePlainAny(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	writePlainAny(nil, levelInfo, "foo")
	assert.Contains(t, buf.String(), "foo")

	buf.Reset()
	writePlainAny(nil, levelDebug, make(chan int))
	assert.Contains(t, buf.String(), "unsupported type")
	writePlainAny(nil, levelDebug, 100)
	assert.Contains(t, buf.String(), "100")

	buf.Reset()
	writePlainAny(nil, levelError, make(chan int))
	assert.Contains(t, buf.String(), "unsupported type")
	writePlainAny(nil, levelSlow, 100)
	assert.Contains(t, buf.String(), "100")

	buf.Reset()
	writePlainAny(hardToWriteWriter{}, levelStat, 100)
	assert.Contains(t, buf.String(), "write error")

	buf.Reset()
	writePlainAny(hardToWriteWriter{}, levelSevere, "foo")
	assert.Contains(t, buf.String(), "write error")

	buf.Reset()
	writePlainAny(hardToWriteWriter{}, levelAlert, "foo")
	assert.Contains(t, buf.String(), "write error")

	buf.Reset()
	writePlainAny(hardToWriteWriter{}, levelFatal, "foo")
	assert.Contains(t, buf.String(), "write error")

	buf.Reset()
	type C struct {
		RC func()
	}
	writePlainAny(nil, levelError, C{
		RC: func() {},
	})
	assert.Contains(t, buf.String(), "runtime/debug.Stack")
}

func TestWritePlainDuplicate(t *testing.T) {
	old := atomic.SwapUint32(&encoding, plainEncodingType)
	t.Cleanup(func() {
		atomic.StoreUint32(&encoding, old)
	})

	var buf bytes.Buffer
	output(&buf, levelInfo, "foo", LogField{
		Key:   "first",
		Value: "a",
	}, LogField{
		Key:   "first",
		Value: "b",
	})
	assert.Contains(t, buf.String(), "foo")
	assert.NotContains(t, buf.String(), "first=a")
	assert.Contains(t, buf.String(), "first=b")

	buf.Reset()
	output(&buf, levelInfo, "foo", LogField{
		Key:   "first",
		Value: "a",
	}, LogField{
		Key:   "first",
		Value: "b",
	}, LogField{
		Key:   "second",
		Value: "c",
	})
	assert.Contains(t, buf.String(), "foo")
	assert.NotContains(t, buf.String(), "first=a")
	assert.Contains(t, buf.String(), "first=b")
	assert.Contains(t, buf.String(), "second=c")
}

func TestLogWithLimitContentLength(t *testing.T) {
	maxLen := atomic.LoadUint32(&maxContentLength)
	atomic.StoreUint32(&maxContentLength, 10)

	t.Cleanup(func() {
		atomic.StoreUint32(&maxContentLength, maxLen)
	})

	t.Run("alert", func(t *testing.T) {
		var buf bytes.Buffer
		w := NewWriter(&buf)
		w.Info("1234567890")
		var v1 mockedEntry
		if err := json.Unmarshal(buf.Bytes(), &v1); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "1234567890", v1.Content)
		assert.False(t, v1.Truncated)

		buf.Reset()
		var v2 mockedEntry
		w.Info("12345678901")
		if err := json.Unmarshal(buf.Bytes(), &v2); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "1234567890", v2.Content)
		assert.True(t, v2.Truncated)
	})
}

func TestComboWriter(t *testing.T) {
	var mockWriters []Writer
	for i := 0; i < 3; i++ {
		mockWriters = append(mockWriters, new(tracedWriter))
	}

	cw := comboWriter{
		writers: mockWriters,
	}

	t.Run("Alert", func(t *testing.T) {
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Alert", "test alert").Once()
		}
		cw.Alert("test alert")
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Alert", "test alert")
		}
	})

	t.Run("Close", func(t *testing.T) {
		for i := range cw.writers {
			if i == 1 {
				cw.writers[i].(*tracedWriter).On("Close").Return(errors.New("error")).Once()
			} else {
				cw.writers[i].(*tracedWriter).On("Close").Return(nil).Once()
			}
		}
		err := cw.Close()
		assert.Error(t, err)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Close")
		}
	})

	t.Run("Debug", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Debug", "test debug", fields).Once()
		}
		cw.Debug("test debug", fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Debug", "test debug", fields)
		}
	})

	t.Run("Error", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Error", "test error", fields).Once()
		}
		cw.Error("test error", fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Error", "test error", fields)
		}
	})

	t.Run("Info", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Info", "test info", fields).Once()
		}
		cw.Info("test info", fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Info", "test info", fields)
		}
	})

	t.Run("Severe", func(t *testing.T) {
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Severe", "test severe").Once()
		}
		cw.Severe("test severe")
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Severe", "test severe")
		}
	})

	t.Run("Slow", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Slow", "test slow", fields).Once()
		}
		cw.Slow("test slow", fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Slow", "test slow", fields)
		}
	})

	t.Run("Stack", func(t *testing.T) {
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Stack", "test stack").Once()
		}
		cw.Stack("test stack")
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Stack", "test stack")
		}
	})

	t.Run("Stat", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Stat", "test stat", fields).Once()
		}
		cw.Stat("test stat", fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Stat", "test stat", fields)
		}
	})
}

type mockedEntry struct {
	Level     string `json:"level"`
	Content   string `json:"content"`
	Truncated bool   `json:"truncated"`
}

type easyToCloseWriter struct{}

func (h easyToCloseWriter) Write(_ []byte) (_ int, _ error) {
	return
}

func (h easyToCloseWriter) Close() error {
	return nil
}

type hardToCloseWriter struct{}

func (h hardToCloseWriter) Write(_ []byte) (_ int, _ error) {
	return
}

func (h hardToCloseWriter) Close() error {
	return errors.New("close error")
}

type hardToWriteWriter struct{}

func (h hardToWriteWriter) Write(_ []byte) (_ int, _ error) {
	return 0, errors.New("write error")
}

type tracedWriter struct {
	mock.Mock
}

func (w *tracedWriter) Alert(v any) {
	w.Called(v)
}

func (w *tracedWriter) Close() error {
	args := w.Called()
	return args.Error(0)
}

func (w *tracedWriter) Debug(v any, fields ...LogField) {
	w.Called(v, fields)
}

func (w *tracedWriter) Error(v any, fields ...LogField) {
	w.Called(v, fields)
}

func (w *tracedWriter) Info(v any, fields ...LogField) {
	w.Called(v, fields)
}

func (w *tracedWriter) Severe(v any) {
	w.Called(v)
}

func (w *tracedWriter) Slow(v any, fields ...LogField) {
	w.Called(v, fields)
}

func (w *tracedWriter) Stack(v any) {
	w.Called(v)
}

func (w *tracedWriter) Stat(v any, fields ...LogField) {
	w.Called(v, fields)
}
