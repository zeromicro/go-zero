package iox

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type (
	textReadOptions struct {
		keepSpace     bool
		withoutBlanks bool
		omitPrefix    string
	}

	// TextReadOption defines the method to customize the text reading functions.
	TextReadOption func(*textReadOptions)
)

// DupReadCloser returns two io.ReadCloser that read from the first will be written to the second.
// The first returned reader needs to be read first, because the content
// read from it will be written to the underlying buffer of the second reader.
func DupReadCloser(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	return io.NopCloser(tee), io.NopCloser(&buf)
}

// DupReadCloserForLargeFile returns two io.ReadCloser that read from the first will be written to the second.
// The first returned reader needs to be read first, because the content
// read from it will be written to the underlying buffer of the second reader.
// the fourth returned close func , After the caller finishes using reader, close the temporary file
func DupReadCloserForLargeFile(reader io.ReadCloser) (io.ReadCloser, io.ReadCloser,error,func()error) {
	f,_:=os.CreateTemp(os.TempDir(),"go-zero-body")
	tee := io.TeeReader(reader, f)
	readTempFile,err:=os.Open(f.Name())
	if err != nil {
		return nil,nil,err,nil
	}
	closeFunc:= func()error {
		return os.Remove(readTempFile.Name())
	}
	return io.NopCloser(tee), io.NopCloser(readTempFile),nil,closeFunc
}

// KeepSpace customizes the reading functions to keep leading and tailing spaces.
func KeepSpace() TextReadOption {
	return func(o *textReadOptions) {
		o.keepSpace = true
	}
}

// LimitDupReadCloser returns two io.ReadCloser that read from the first will be written to the second.
// But the second io.ReadCloser is limited to up to n bytes.
// The first returned reader needs to be read first, because the content
// read from it will be written to the underlying buffer of the second reader.
func LimitDupReadCloser(reader io.ReadCloser, n int64) (io.ReadCloser, io.ReadCloser) {
	var buf bytes.Buffer
	tee := LimitTeeReader(reader, &buf, n)
	return io.NopCloser(tee), io.NopCloser(&buf)
}

// ReadBytes reads exactly the bytes with the length of len(buf)
func ReadBytes(reader io.Reader, buf []byte) error {
	var got int

	for got < len(buf) {
		n, err := reader.Read(buf[got:])
		if err != nil {
			return err
		}

		got += n
	}

	return nil
}

// ReadText reads content from the given file with leading and tailing spaces trimmed.
func ReadText(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// ReadTextLines reads the text lines from given file.
func ReadTextLines(filename string, opts ...TextReadOption) ([]string, error) {
	var readOpts textReadOptions
	for _, opt := range opts {
		opt(&readOpts)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !readOpts.keepSpace {
			line = strings.TrimSpace(line)
		}
		if readOpts.withoutBlanks && len(line) == 0 {
			continue
		}
		if len(readOpts.omitPrefix) > 0 && strings.HasPrefix(line, readOpts.omitPrefix) {
			continue
		}

		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

// WithoutBlank customizes the reading functions to ignore blank lines.
func WithoutBlank() TextReadOption {
	return func(o *textReadOptions) {
		o.withoutBlanks = true
	}
}

// OmitWithPrefix customizes the reading functions to ignore the lines with given leading prefix.
func OmitWithPrefix(prefix string) TextReadOption {
	return func(o *textReadOptions) {
		o.omitPrefix = prefix
	}
}
