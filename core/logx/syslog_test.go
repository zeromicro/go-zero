package logx

import (
	"encoding/json"
	"log"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testlog = "Stay hungry, stay foolish."

func TestCollectSysLog(t *testing.T) {
	CollectSysLog()
	content := getContent(captureOutput(func() {
		log.Print(testlog)
	}))
	assert.True(t, strings.Contains(content, testlog))
}

func TestRedirector(t *testing.T) {
	var r redirector
	content := getContent(captureOutput(func() {
		r.Write([]byte(testlog))
	}))
	assert.Equal(t, testlog, content)
}

func captureOutput(f func()) string {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	prevLevel := atomic.LoadUint32(&logLevel)
	SetLevel(InfoLevel)
	f()
	SetLevel(prevLevel)

	return w.String()
}

func getContent(jsonStr string) string {
	var entry logEntry
	json.Unmarshal([]byte(jsonStr), &entry)
	val, ok := entry.Content.(string)
	if ok {
		return val
	}
	return ""
}
