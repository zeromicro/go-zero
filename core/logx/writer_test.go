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
	w.Info(literal, 0)
	assert.Contains(t, buf.String(), literal)
	buf.Reset()
	w.Debug(literal, 0)
	assert.Contains(t, buf.String(), literal)
}

func TestConsoleWriter(t *testing.T) {
	var buf bytes.Buffer
	w := newConsoleWriter()
	lw := newLogWriter(log.New(&buf, "", 0))
	w.(*concreteWriter).errorLog = lw
	w.Alert("foo bar 1", 0)
	var val mockedEntry
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelAlert, val.Level)
	assert.Equal(t, "foo bar 1", val.Content)

	buf.Reset()
	w.(*concreteWriter).errorLog = lw
	w.Error("foo bar 2", 0)
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelError, val.Level)
	assert.Equal(t, "foo bar 2", val.Content)

	buf.Reset()
	w.(*concreteWriter).infoLog = lw
	w.Info("foo bar 3", 0)
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelInfo, val.Level)
	assert.Equal(t, "foo bar 3", val.Content)

	buf.Reset()
	w.(*concreteWriter).severeLog = lw
	w.Severe("foo bar 4", 0)
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelFatal, val.Level)
	assert.Equal(t, "foo bar 4", val.Content)

	buf.Reset()
	w.(*concreteWriter).slowLog = lw
	w.Slow("foo bar 5", 0)
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, levelSlow, val.Level)
	assert.Equal(t, "foo bar 5", val.Content)

	buf.Reset()
	w.(*concreteWriter).statLog = lw
	w.Stat("foo bar 6", 0)
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
		w.Alert("foo", 0)
		w.Debug("foo", 0)
		w.Error("foo", 0)
		w.Info("foo", 0)
		w.Severe("foo", 0)
		w.Stack("foo", 0)
		w.Stat("foo", 0)
		w.Slow("foo", 0)
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
	output(&buf, levelInfo, "foo", 0, LogField{
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
	output(&buf, levelInfo, "foo", 0, LogField{
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

func TestLogWithSensitive(t *testing.T) {
	old := atomic.SwapUint32(&encoding, plainEncodingType)
	t.Cleanup(func() {
		atomic.StoreUint32(&encoding, old)
	})

	t.Run("sensitive", func(t *testing.T) {
		var buf bytes.Buffer
		output(&buf, levelInfo, User{
			Name: "kevin",
			Pass: "123",
		}, 0, LogField{
			Key:   "first",
			Value: "a",
		}, LogField{
			Key:   "first",
			Value: "b",
		})
		assert.Contains(t, buf.String(), maskedContent)
		assert.NotContains(t, buf.String(), "first=a")
		assert.Contains(t, buf.String(), "first=b")
	})

	t.Run("sensitive fields", func(t *testing.T) {
		var buf bytes.Buffer
		output(&buf, levelInfo, "foo", 0, LogField{
			Key: "first",
			Value: User{
				Name: "kevin",
				Pass: "123",
			},
		}, LogField{
			Key:   "second",
			Value: "b",
		})
		assert.Contains(t, buf.String(), "foo")
		assert.Contains(t, buf.String(), "first")
		assert.Contains(t, buf.String(), maskedContent)
		assert.Contains(t, buf.String(), "second=b")
	})
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
		w.Info("1234567890", 0)
		var v1 mockedEntry
		if err := json.Unmarshal(buf.Bytes(), &v1); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "1234567890", v1.Content)
		assert.False(t, v1.Truncated)

		buf.Reset()
		var v2 mockedEntry
		w.Info("12345678901", 0)
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
			mw.(*tracedWriter).On("Alert", "test alert", uint64(0)).Once()
		}
		cw.Alert("test alert", 0)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Alert", "test alert", uint64(0))
		}
	})

	t.Run("Debug", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Debug", "test debug", uint64(0), fields).Once()
		}
		cw.Debug("test debug", 0, fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Debug", "test debug", uint64(0), fields)
		}
	})

	t.Run("Error", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Error", "test error", uint64(0), fields).Once()
		}
		cw.Error("test error", 0, fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Error", "test error", uint64(0), fields)
		}
	})

	t.Run("Info", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Info", "test info", uint64(0), fields).Once()
		}
		cw.Info("test info", 0, fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Info", "test info", uint64(0), fields)
		}
	})

	t.Run("Severe", func(t *testing.T) {
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Severe", "test severe", uint64(0)).Once()
		}
		cw.Severe("test severe", 0)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Severe", "test severe", uint64(0))
		}
	})

	t.Run("Slow", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Slow", "test slow", uint64(0), fields).Once()
		}
		cw.Slow("test slow", 0, fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Slow", "test slow", uint64(0), fields)
		}
	})

	t.Run("Stack", func(t *testing.T) {
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Stack", "test stack", uint64(0)).Once()
		}
		cw.Stack("test stack", 0)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Stack", "test stack", uint64(0))
		}
	})

	t.Run("Stat", func(t *testing.T) {
		fields := []LogField{{Key: "key", Value: "value"}}
		for _, mw := range cw.writers {
			mw.(*tracedWriter).On("Stat", "test stat", uint64(0), fields).Once()
		}
		cw.Stat("test stat", 0, fields...)
		for _, mw := range cw.writers {
			mw.(*tracedWriter).AssertCalled(t, "Stat", "test stat", uint64(0), fields)
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

func (w *tracedWriter) Alert(v any, loggerID uint64) {
	w.Called(v, loggerID)
}

func (w *tracedWriter) Close() error {
	args := w.Called()
	return args.Error(0)
}

func (w *tracedWriter) Debug(v any, loggerID uint64, fields ...LogField) {
	w.Called(v, loggerID, fields)
}

func (w *tracedWriter) Error(v any, loggerID uint64, fields ...LogField) {
	w.Called(v, loggerID, fields)
}

func (w *tracedWriter) Info(v any, loggerID uint64, fields ...LogField) {
	w.Called(v, loggerID, fields)
}

func (w *tracedWriter) Severe(v any, loggerID uint64) {
	w.Called(v, loggerID)
}

func (w *tracedWriter) Slow(v any, loggerID uint64, fields ...LogField) {
	w.Called(v, loggerID, fields)
}

func (w *tracedWriter) Stack(v any, loggerID uint64) {
	w.Called(v, loggerID)
}

func (w *tracedWriter) Stat(v any, loggerID uint64, fields ...LogField) {
	w.Called(v, loggerID, fields)
}
