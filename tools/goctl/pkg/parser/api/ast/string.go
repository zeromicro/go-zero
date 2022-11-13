package ast

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type Writer struct {
	tw             *tabwriter.Writer
	lastWriteToken token.Token
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		tw:             tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		lastWriteToken: token.NewIllegalToken(0, token.IllegalPosition),
	}
}

func (w *Writer) Write(prefix string, toks ...token.Token) {
	if len(toks) == 0 {
		return
	}

	var lastWriteToken = w.lastWriteToken
	var inOneLine = true
	var one = toks[0]
	var list []string
	for _, tok := range toks {
		if one.Line() != tok.Line() {
			inOneLine = false
		}
		list = append(list, tok.Text)
		w.lastWriteToken = tok
	}
	if inOneLine {
		if one.Line() > w.lastWriteToken.Line() {
			w.NewLine()
		}
		_, _ = fmt.Fprint(w.tw, strings.Join(list, "\t"))
		return
	}

	w.lastWriteToken = lastWriteToken
	_, _ = fmt.Fprint(w.tw, prefix)
	for idx, tok := range toks {
		if tok.Line() > w.lastWriteToken.Line() {
			w.NewLine()
		}
		_, _ = fmt.Fprint(w.tw, tok.Text)
		if idx < len(toks)-1 {
			_, _ = fmt.Fprint(w.tw, " ")
		}
		w.lastWriteToken = tok
	}
}

func (w *Writer) NewLine() {
	_, _ = fmt.Fprint(w.tw, "\n")
}
