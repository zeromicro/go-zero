package filex

import (
	"errors"
	"os"
)

// errExceedFileSize indicates that the file size is exceeded.
var errExceedFileSize = errors.New("exceed file size")

// A RangeReader is used to read a range of content from a file.
type RangeReader struct {
	file  *os.File
	start int64
	stop  int64
}

// NewRangeReader returns a RangeReader, which will read the range of content from file.
func NewRangeReader(file *os.File, start, stop int64) *RangeReader {
	return &RangeReader{
		file:  file,
		start: start,
		stop:  stop,
	}
}

// Read reads the range of content into p.
func (rr *RangeReader) Read(p []byte) (n int, err error) {
	stat, err := rr.file.Stat()
	if err != nil {
		return 0, err
	}

	if rr.stop < rr.start || rr.start >= stat.Size() {
		return 0, errExceedFileSize
	}

	if rr.stop-rr.start < int64(len(p)) {
		p = p[:rr.stop-rr.start]
	}

	n, err = rr.file.ReadAt(p, rr.start)
	if err != nil {
		return n, err
	}

	rr.start += int64(n)
	return
}
