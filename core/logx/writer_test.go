package logx

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWriter(t *testing.T) {
	const literal = "foo bar"
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.Info(literal)
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

func TestNopWriter(t *testing.T) {
	assert.NotPanics(t, func() {
		var w nopWriter
		w.Alert("foo")
		w.Error("foo")
		w.Info("foo")
		w.Severe("foo")
		w.Stack("foo")
		w.Stat("foo")
		w.Slow("foo")
		w.Close()
	})
}

func TestWriteJson(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	writeJson(nil, "foo")
	assert.Contains(t, buf.String(), "foo")
	buf.Reset()
	writeJson(nil, make(chan int))
	assert.Contains(t, buf.String(), "unsupported type")
}

func TestWritePlainAny(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	writePlainAny(nil, levelInfo, "foo")
	assert.Contains(t, buf.String(), "foo")

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

}

type mockedEntry struct {
	Level   string `json:"level"`
	Content string `json:"content"`
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
