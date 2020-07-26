package filex

import (
	"errors"
	"os"
)

type RangeReader struct {
	file  *os.File
	start int64
	stop  int64
}

func NewRangeReader(file *os.File, start, stop int64) *RangeReader {
	return &RangeReader{
		file:  file,
		start: start,
		stop:  stop,
	}
}

func (rr *RangeReader) Read(p []byte) (n int, err error) {
	stat, err := rr.file.Stat()
	if err != nil {
		return 0, err
	}

	if rr.stop < rr.start || rr.start >= stat.Size() {
		return 0, errors.New("exceed file size")
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
