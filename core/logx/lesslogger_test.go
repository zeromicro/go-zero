package logx

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessLogger_Error(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	l := NewLessLogger(500)
	for i := 0; i < 100; i++ {
		l.Error("hello")
	}

	assert.Equal(t, 1, strings.Count(builder.String(), "\n"))
}

func TestLessLogger_Errorf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	l := NewLessLogger(500)
	for i := 0; i < 100; i++ {
		l.Errorf("hello")
	}

	assert.Equal(t, 1, strings.Count(builder.String(), "\n"))
}
