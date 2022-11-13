package format

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

const nilIndent = ""
const whiteSpace = " "
const indent = "\t"
const newLine = "\n"

type Formatter interface {
	Write(prefix string, toks ...token.Token)
	WriteInOneLine(prefix string, toks ...token.Token)
	WriteBetween(prefix string, left, right token.Token)
	WriteSpaceInfix(prefix string, toks ...token.Token)
	WriteSpaceInfixBetween(prefix string, left, right token.Token)
	Skip(left, right token.Token)
	NewLine()
	Flush()
}

type NilWriter struct{}

func NewNilWriter() *NilWriter {
	return &NilWriter{}
}
func (n NilWriter) Write(string, ...token.Token)                            {}
func (n NilWriter) WriteInOneLine(string, ...token.Token)                   {}
func (n NilWriter) WriteSpaceInfix(string, ...token.Token)                  {}
func (n NilWriter) WriteBetween(string, token.Token, token.Token)           {}
func (n NilWriter) WriteSpaceInfixBetween(string, token.Token, token.Token) {}
func (n NilWriter) Skip(token.Token, token.Token)                           {}
func (n NilWriter) NewLine()                                                {}
func (n NilWriter) Flush()                                                  {}

type Writer struct {
	tw             *tabwriter.Writer
	lastWriteToken token.Token
	tokenSet       *token.Set
}

func NewWriter(writer io.Writer, tokenSet *token.Set) *Writer {
	return &Writer{
		tw:             tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		lastWriteToken: token.InitToken,
		tokenSet:       tokenSet,
	}
}

func (w *Writer) WriteBetween(prefix string, left, right token.Token) {
	tokens := w.tokenSet.Between(left, right, token.AllIn)
	if len(tokens) > 0 {
		w.Write(prefix, tokens...)
	}
}

func (w *Writer) WriteSpaceInfixBetween(prefix string, left, right token.Token) {
	tokens := w.tokenSet.Between(left, right, token.AllIn)
	if len(tokens) > 0 {
		w.WriteSpaceInfix(prefix, tokens...)
	}
}

func (w *Writer) WriteSpaceInfix(prefix string, toks ...token.Token) {
	if len(toks) == 0 {
		return
	}

	defer func() {
		tail := toks[len(toks)-1]
		lineAfter := w.tokenSet.LineCommentAfter(tail)
		if len(lineAfter) > 0 {
			w.write(nilIndent, lineAfter...)
		}
	}()
	var hasDoc = false
	for _, e := range toks {
		if e.IsDocument()||e.IsComment() {
			hasDoc = true
			break
		}
	}

	var one = toks[0]
	gaps := w.tokenSet.Between(w.lastWriteToken, one, token.NotIn)
	if len(gaps) > 0 {
		w.write(nilIndent, gaps...)
	}

	_, _ = fmt.Fprint(w.tw, prefix)
	for idx, tok := range toks {
		if tok.Line() > w.lastWriteToken.Line() && (idx == 0 || hasDoc) {
			w.newLine()
		}
		_, _ = fmt.Fprint(w.tw, tok.Text)
		if idx < len(toks)-1 {
			_, _ = fmt.Fprint(w.tw, whiteSpace)
		}
		w.lastWriteToken = tok
	}
}

func (w *Writer) WriteInOneLine(prefix string, toks ...token.Token) {
	if len(toks) == 0 {
		return
	}
	defer func() {
		tail := toks[len(toks)-1]
		lineAfter := w.tokenSet.LineCommentAfter(tail)
		if len(lineAfter) > 0 {
			w.write(nilIndent, lineAfter...)
		}
	}()
	var one = toks[0]
	gaps := w.tokenSet.Between(w.lastWriteToken, one, token.NotIn)
	if len(gaps) > 0 {
		w.write(nilIndent, gaps...)
	}

	var lastWriteToken = w.lastWriteToken
	var hasDocument = false
	var list []string
	for _, tok := range toks {
		if tok.IsDocument() {
			hasDocument = true
		}
		list = append(list, tok.Text)
		lastWriteToken = tok
	}
	if !hasDocument {
		_, _ = fmt.Fprint(w.tw, prefix)
		_, _ = fmt.Fprint(w.tw, strings.Join(list, whiteSpace))
		w.lastWriteToken = lastWriteToken
		return
	}

	w.write(prefix, toks...)
}

func (w *Writer) Write(prefix string, toks ...token.Token) {
	if len(toks) == 0 {
		return
	}

	defer func() {
		tail := toks[len(toks)-1]
		lineAfter := w.tokenSet.LineCommentAfter(tail)
		if len(lineAfter) > 0 {
			w.write(nilIndent, lineAfter...)
		}
	}()
	var one = toks[0]
	gaps := w.tokenSet.Between(w.lastWriteToken, one, token.NotIn)
	if len(gaps) > 0 {
		w.write(nilIndent, gaps...)
	}

	var lastWriteToken = w.lastWriteToken
	var inOneLine = true
	var list []string
	for _, tok := range toks {
		if one.Line() != tok.Line() {
			inOneLine = false
		}
		list = append(list, tok.Text)
		lastWriteToken = tok
	}
	if inOneLine {
		if one.Line() > w.lastWriteToken.Line() {
			w.newLine()
		}
		_, _ = fmt.Fprint(w.tw, prefix)
		_, _ = fmt.Fprint(w.tw, strings.Join(list, indent))
		w.lastWriteToken = lastWriteToken
		return
	}

	w.write(prefix, toks...)
}

func (w *Writer) write(prefix string, toks ...token.Token) {
	_, _ = fmt.Fprint(w.tw, prefix)
	for idx, tok := range toks {
		if tok.Line() > w.lastWriteToken.Line() {
			w.newLine()
		}
		_, _ = fmt.Fprint(w.tw, tok.Text)
		if idx < len(toks)-1 {
			_, _ = fmt.Fprint(w.tw, whiteSpace)
		}
		w.lastWriteToken = tok
	}
}

func (w *Writer) Skip(left, right token.Token) {
	gaps := w.tokenSet.Between(w.lastWriteToken, left, token.NotIn)
	if len(gaps) > 0 {
		w.write("", gaps...)
	}
	w.lastWriteToken = right
}

func (w *Writer) NewLine() {
	_, _ = fmt.Fprint(w.tw, newLine)
}

func (w *Writer) Flush() {
	list := w.tokenSet.Between(w.lastWriteToken, w.tokenSet.LastToken(), token.RightIn)
	w.Write(nilIndent, list...)
	w.NewLine()
	_ = w.tw.Flush()
}

func (w *Writer) newLine() {
	_, _ = fmt.Fprint(w.tw, newLine)
}
