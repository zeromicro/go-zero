package iox

import "io"

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

// NopCloser returns an io.WriteCloser that does nothing on calling Close.
func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}
