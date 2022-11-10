package ast

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

var infixWhiteSpace = withInfix(whiteSpace)
var infixIndent = withInfix(indent)

type infixImpl struct {
	infix string
}

func (i infixImpl) Infix() string {
	return i.infix
}

func withInfix(infix string) infixImpl {
	return infixImpl{infix: infix}
}

type IInfix interface {
	Infix() string
}

type Writer struct {
	w *bytes.Buffer
}

type TabWriter struct {
	tw *tabwriter.Writer
}

func NewTabWriter(w *bytes.Buffer) *TabWriter {
	return &TabWriter{
		tw: tabwriter.NewWriter(w, 1, 8, 1, ' ', tabwriter.TabIndent),
	}
}

func (tw *TabWriter) WriteWithInfixIndentln(prefix string, v ...interface{}) {
	_, _ = fmt.Fprint(tw.tw, prefix)
	_, _ = fmt.Fprint(tw.tw, sprint(append([]interface{}{infixIndent}, v...)...), "\n")
}

func (tw *TabWriter) Flush() {
	_ = tw.tw.Flush()
}

func NewWriter() *Writer {
	w := bytes.NewBuffer(nil)
	return &Writer{
		w: w,
	}
}

func (tw *Writer) UseTabWriter() *TabWriter {
	return NewTabWriter(tw.w)
}

func (tw *Writer) Write(v ...interface{}) {
	tw.w.WriteString(sprint(v...))
}

func (tw *Writer) WriteWithWhiteSpaceInfix(prefix string, v ...interface{}) {
	tw.Write(prefix, sprint(append([]interface{}{infixWhiteSpace}, v...)...))
}

func (tw *Writer) WriteWithWhiteSpaceInfixln(prefix string, v ...interface{}) {
	tw.WriteWithWhiteSpaceInfix(prefix, v...)
	tw.NewLine()
}

func (tw *Writer) Writeln(v ...interface{}) {
	tw.Write(v...)
	tw.NewLine()
}

func (tw *Writer) NewLine() {
	tw.w.WriteByte('\n')
}

func (tw *Writer) String() string {
	return tw.w.String()
}

func sprint(v ...interface{}) string {
	var data []string
	var infix string
	for _, val := range v {
		switch elem := val.(type) {
		case IInfix:
			infix = elem.Infix()
		default:
			value := getString(val)
			if len(value) == 0 {
				continue
			}
			data = append(data, value)
		}
	}
	return strings.Join(data, infix)
}

func getString(v interface{}) string {
	switch val := v.(type) {
	case token.Token:
		return val.Text
	case fmt.Stringer:
		return val.String()
	case IInfix:
		return val.Infix()
	case string:
		return val
	default:
		return fmt.Sprint(v)
	}
}
