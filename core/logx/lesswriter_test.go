package logx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLessWriter(t *testing.T) {
	var builder strings.Builder
	w := NewLessWriter(&builder, 500)
	for i := 0; i < 100; i++ {
		_, err := w.Write([]byte("hello"))
		assert.Nil(t, err)
	}

	assert.Equal(t, "hello", builder.String())
}
