package iox

import (
	"bufio"
	"io"
	"strings"
)

type TextLineScanner struct {
	reader  *bufio.Reader
	hasNext bool
	line    string
	err     error
}

func NewTextLineScanner(reader io.Reader) *TextLineScanner {
	return &TextLineScanner{
		reader:  bufio.NewReader(reader),
		hasNext: true,
	}
}

func (scanner *TextLineScanner) Scan() bool {
	if !scanner.hasNext {
		return false
	}

	line, err := scanner.reader.ReadString('\n')
	scanner.line = strings.TrimRight(line, "\n")
	if err == io.EOF {
		scanner.hasNext = false
		return true
	} else if err != nil {
		scanner.err = err
		return false
	}
	return true
}

func (scanner *TextLineScanner) Line() (string, error) {
	return scanner.line, scanner.err
}
