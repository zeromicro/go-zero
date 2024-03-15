package iox

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// A TextLineScanner is a scanner that can scan lines from the given reader.
type TextLineScanner struct {
	reader  *bufio.Reader
	hasNext bool
	line    string
	err     error
}

// NewTextLineScanner returns a TextLineScanner with the given reader.
func NewTextLineScanner(reader io.Reader) *TextLineScanner {
	return &TextLineScanner{
		reader:  bufio.NewReader(reader),
		hasNext: true,
	}
}

// Scan checks if scanner has more lines to read.
func (scanner *TextLineScanner) Scan() bool {
	if !scanner.hasNext {
		return false
	}

	line, err := scanner.reader.ReadString('\n')
	scanner.line = strings.TrimRight(line, "\n")
	if errors.Is(err, io.EOF) {
		scanner.hasNext = false
		return true
	} else if err != nil {
		scanner.err = err
		return false
	}
	return true
}

// Line returns the next available line.
func (scanner *TextLineScanner) Line() (string, error) {
	return scanner.line, scanner.err
}
