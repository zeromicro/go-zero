package stat

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

func TestMetrics(t *testing.T) {
	DisableLog()
	defer logEnabled.Set(true)

	counts := []int{1, 5, 10, 100, 1000, 1000}
	for _, count := range counts {
		m := NewMetrics("foo")
		m.SetName("bar")
		for i := 0; i < count; i++ {
			m.Add(Task{
				Duration:    time.Millisecond * time.Duration(i),
				Description: strconv.Itoa(i),
			})
		}
		m.AddDrop()
		var writer mockedWriter
		SetReportWriter(&writer)
		m.executor.Flush()
		assert.Equal(t, "bar", writer.report.Name)
	}
}

func TestTopDurationWithEmpty(t *testing.T) {
	assert.Equal(t, float32(0), getTopDuration(nil))
	assert.Equal(t, float32(0), getTopDuration([]Task{}))
}

func TestLogAndReport(t *testing.T) {
	buf := logtest.NewCollector(t)
	old := logEnabled.True()
	logEnabled.Set(true)
	t.Cleanup(func() {
		logEnabled.Set(old)
	})

	log(&StatReport{})
	assert.NotEmpty(t, buf.String())

	writerLock.Lock()
	writer := reportWriter
	writerLock.Unlock()
	buf = logtest.NewCollector(t)
	t.Cleanup(func() {
		SetReportWriter(writer)
	})
	SetReportWriter(&badWriter{})
	writeReport(&StatReport{})
	assert.NotEmpty(t, buf.String())
}

type mockedWriter struct {
	report *StatReport
}

func (m *mockedWriter) Write(report *StatReport) error {
	m.report = report
	return nil
}

type badWriter struct{}

func (b *badWriter) Write(_ *StatReport) error {
	return errors.New("bad")
}
