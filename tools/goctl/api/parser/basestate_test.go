package parser

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProperties(t *testing.T) {
	const text = `(summary: hello world)`
	var builder bytes.Buffer
	builder.WriteString(text)
	var lineNumber = 1
	var state = newBaseState(bufio.NewReader(&builder), &lineNumber)
	m, err := state.parseProperties()
	assert.Nil(t, err)
	assert.Equal(t, "hello world", m["summary"])
}
