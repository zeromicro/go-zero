package iox

import (
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	const val = `1
2
3
4`
	reader := strings.NewReader(val)
	scanner := NewTextLineScanner(reader)
	var lines []string
	for scanner.Scan() {
		line, err := scanner.Line()
		assert.Nil(t, err)
		lines = append(lines, line)
	}
	assert.EqualValues(t, []string{"1", "2", "3", "4"}, lines)
}

func TestScannerDoesNotReturnEmptyEOF(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: ""},
		{input: "one\n", want: []string{"one"}},
		{input: "\n", want: []string{""}},
		{input: "one\n\n", want: []string{"one", ""}},
	}

	for _, test := range tests {
		scanner := NewTextLineScanner(strings.NewReader(test.input))
		var lines []string
		for scanner.Scan() {
			line, err := scanner.Line()
			assert.NoError(t, err)
			lines = append(lines, line)
		}
		assert.Equal(t, test.want, lines)
	}
}

func TestBadScanner(t *testing.T) {
	scanner := NewTextLineScanner(iotest.ErrReader(iotest.ErrTimeout))
	assert.False(t, scanner.Scan())
	_, err := scanner.Line()
	assert.ErrorIs(t, err, iotest.ErrTimeout)
}
