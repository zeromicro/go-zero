package logx

import "io"

type lessWriter struct {
	*limitedExecutor
	writer io.Writer
}

func newLessWriter(writer io.Writer, milliseconds int) *lessWriter {
	return &lessWriter{
		limitedExecutor: newLimitedExecutor(milliseconds),
		writer:          writer,
	}
}

func (w *lessWriter) Write(p []byte) (n int, err error) {
	w.logOrDiscard(func() {
		w.writer.Write(p)
	})
	return len(p), nil
}
