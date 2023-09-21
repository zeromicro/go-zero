package iox

import "io"

// LimitTeeReader returns a Reader that writes up to n bytes to w what it reads from r.
// First n bytes reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the first n bytes read completes.
// Any error encountered while writing is reported as a read error.
func LimitTeeReader(r io.Reader, w io.Writer, n int64) io.Reader {
	return &limitTeeReader{r, w, n}
}

type limitTeeReader struct {
	r io.Reader
	w io.Writer
	n int64 // limit bytes remaining
}

func (t *limitTeeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 && t.n > 0 {
		limit := int64(n)
		if limit > t.n {
			limit = t.n
		}
		if n, err := t.w.Write(p[:limit]); err != nil {
			return n, err
		}

		t.n -= limit
	}

	return
}
