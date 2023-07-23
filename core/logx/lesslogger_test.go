package logx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessLogger_Error(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	l := NewLessLogger(500)
	for i := 0; i < 100; i++ {
		l.Error("hello")
	}

	assert.Equal(t, 1, strings.Count(w.String(), "\n"))
}

func TestLessLogger_Errorf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	l := NewLessLogger(500)
	for i := 0; i < 100; i++ {
		l.Errorf("hello")
	}

	assert.Equal(t, 1, strings.Count(w.String(), "\n"))
}
