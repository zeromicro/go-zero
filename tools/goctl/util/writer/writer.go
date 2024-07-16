package writer

import (
	"bytes"
	"fmt"
	"strings"
)

type Writer struct {
	w      *bytes.Buffer
	indent string
}

func New(indent string) *Writer {
	return &Writer{
		w:      bytes.NewBuffer(nil),
		indent: indent,
	}
}

func (w *Writer) WriteStringln(s string) {
	w.w.WriteString(s)
	w.NewLine()
}

func (w *Writer) WriteWithIndentStringln(s string) {
	w.w.WriteString(w.indent)
	w.WriteStringln(s)
}

func (w *Writer) WriteWithIndentStringf(format string, a ...any) {
	w.w.WriteString(w.indent)
	w.Writef(format, a...)
}

func (w *Writer) Writef(format string, a ...any) {
	w.w.WriteString(fmt.Sprintf(format, a...))
}

func (w *Writer) NewLine() {
	w.w.WriteRune('\n')
}

func (w *Writer) UndoNewLine() {
	w.Undo("\n")
}

func (w *Writer) Undo(s string) {
	val := w.w.String()
	w.w.Reset()
	w.w.WriteString(strings.TrimSuffix(val, s))
}

func (w *Writer) Remove(s string) {
	val := w.w.String()
	w.w.Reset()
	val = strings.ReplaceAll(val, s, "")
	w.w.WriteString(val)
}

func (w *Writer) String() string {
	return w.w.String()
}

func (w *Writer) Bytes() []byte {
	return w.w.Bytes()
}
