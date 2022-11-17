package ast

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

const (
	NilIndent  = ""
	WhiteSpace = " "
	Indent     = "\t"
	NewLine    = "\n"
)

const (
	_ WriteMode = 1 << iota
	ModeAuto
	ModeExpectInSameLine
)

type Option func(o *option)

type option struct {
	prefix  string
	infix   string
	mode    WriteMode
	nodes   []Node
	rawText bool
}

type WriteMode int

type Writer struct {
	tw     *tabwriter.Writer
	writer io.Writer
}

func WithNode(nodes ...Node) Option {
	return func(o *option) {
		o.nodes = nodes
	}
}

func WithMode(mode WriteMode) Option {
	return func(o *option) {
		o.mode = mode
	}
}

func WithPrefix(prefix ...string) Option {
	return func(o *option) {
		for _, p := range prefix {
			o.prefix = p
		}
	}
}

func WithInfix(infix string) Option {
	return func(o *option) {
		o.infix = infix
	}
}

func WithRawText() Option {
	return func(o *option) {
		o.rawText = true
	}
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		tw:     tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		writer: writer,
	}
}

func NewBufferWriter() *Writer {
	writer := bytes.NewBuffer(nil)
	return &Writer{
		tw:     tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		writer: writer,
	}
}

func (w *Writer) String() string {
	buffer, ok := w.writer.(*bytes.Buffer)
	if !ok {
		return ""
	}
	w.Flush()
	return buffer.String()
}

func (w *Writer) Flush() {
	_ = w.tw.Flush()
}

func (w *Writer) NewLine() {
	_, _ = fmt.Fprint(w.tw, NewLine)
}

func (w *Writer) Write(opts ...Option) {
	if len(opts) == 0 {
		return
	}

	var opt = new(option)
	opt.mode = ModeAuto
	opt.prefix = NilIndent
	opt.infix = WhiteSpace
	for _, v := range opts {
		v(opt)
	}

	w.write(opt)
}

func (w *Writer) WriteText(text string) {
	_, _ = fmt.Fprintf(w.tw, text)
}

func (w *Writer) write(opt *option) {
	if len(opt.nodes) == 0 {
		return
	}

	var textList []string
	line := opt.nodes[0].End().Line
	for idx, node := range opt.nodes {
		tokenNode, ok := node.(*TokenNode)
		if ok && tokenNode.HasLeadingCommentGroup() && idx < len(opt.nodes)-1 {
			opt.mode = ModeAuto
		}
		if opt.mode == ModeAuto && node.Pos().Line > line {
			textList = append(textList, NewLine)
		}
		line = node.End().Line
		if util.TrimWhiteSpace(node.Format(opt.prefix)) == "" {
			continue
		}

		textList = append(textList, node.Format(opt.prefix))
	}

	if opt.rawText {
		_, _ = fmt.Fprint(w.writer, strings.Join(textList, opt.infix))
		return
	}
	_, _ = fmt.Fprint(w.tw, strings.Join(textList, opt.infix))
}
