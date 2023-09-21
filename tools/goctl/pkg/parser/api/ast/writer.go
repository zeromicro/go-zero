package ast

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
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
	// ModeAuto is the default mode, which will automatically
	//determine whether to write a newline.
	ModeAuto

	// ModeExpectInSameLine will write in the same line.
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

type tokenNodeOption func(o *tokenNodeOpt)
type tokenNodeOpt struct {
	prefix               string
	infix                string
	ignoreHeadComment    bool
	ignoreLeadingComment bool
}

// WriteMode is the mode of writing.
type WriteMode int

// Writer is the writer of ast.
type Writer struct {
	tw     *tabwriter.Writer
	writer io.Writer
}

func transfer2TokenNode(node Node, isChild bool, opt ...tokenNodeOption) *TokenNode {
	option := new(tokenNodeOpt)
	for _, o := range opt {
		o(option)
	}

	var copyOpt = append([]tokenNodeOption(nil), opt...)
	var tn *TokenNode
	switch val := node.(type) {
	case *AnyDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.Any, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.Any = tn
	case *ArrayDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.LBrack, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.LBrack = tn
	case *BaseDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.Base, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.Base = tn
	case *InterfaceDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.Interface, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.Interface = tn
	case *MapDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.Map, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.Map = tn
	case *PointerDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.Star, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.Star = tn
	case *SliceDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.LBrack, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.LBrack = tn
	case *StructDataType:
		copyOpt = append(copyOpt, withTokenNodePrefix(NilIndent))
		tn = transferTokenNode(val.LBrace, copyOpt...)
		if option.ignoreHeadComment {
			tn.HeadCommentGroup = nil
		}
		if option.ignoreLeadingComment {
			tn.LeadingCommentGroup = nil
		}
		val.isChild = isChild
		val.LBrace = tn
	default:
	}

	return &TokenNode{
		headFlag:    node.HasHeadCommentGroup(),
		leadingFlag: node.HasLeadingCommentGroup(),
		Token: token.Token{
			Text:     node.Format(option.prefix),
			Position: node.Pos(),
		},
		LeadingCommentGroup: CommentGroup{
			{
				token.Token{Position: node.End()},
			},
		},
	}
}

func transferNilInfixNode(nodes []*TokenNode, opt ...tokenNodeOption) *TokenNode {
	result := &TokenNode{}
	var option = new(tokenNodeOpt)
	for _, o := range opt {
		o(option)
	}

	var list []string
	for _, n := range nodes {
		list = append(list, n.Token.Text)
	}

	result.Token = token.Token{
		Text:     option.prefix + strings.Join(list, option.infix),
		Position: nodes[0].Pos(),
	}

	if !option.ignoreHeadComment {
		result.HeadCommentGroup = nodes[0].HeadCommentGroup
	}
	if !option.ignoreLeadingComment {
		result.LeadingCommentGroup = nodes[len(nodes)-1].LeadingCommentGroup
	}

	return result
}

func transferTokenNode(node *TokenNode, opt ...tokenNodeOption) *TokenNode {
	result := &TokenNode{}
	var option = new(tokenNodeOpt)
	for _, o := range opt {
		o(option)
	}
	result.Token = token.Token{
		Type:     node.Token.Type,
		Text:     option.prefix + node.Token.Text,
		Position: node.Token.Position,
	}
	if !option.ignoreHeadComment {
		for _, v := range node.HeadCommentGroup {
			result.HeadCommentGroup = append(result.HeadCommentGroup,
				&CommentStmt{Comment: token.Token{
					Type:     v.Comment.Type,
					Text:     option.prefix + v.Comment.Text,
					Position: v.Comment.Position,
				}})
		}
	}
	if !option.ignoreLeadingComment {
		result.LeadingCommentGroup = append(result.LeadingCommentGroup, node.LeadingCommentGroup...)
	}
	return result
}

func ignoreHeadComment() tokenNodeOption {
	return func(o *tokenNodeOpt) {
		o.ignoreHeadComment = true
	}
}

func ignoreLeadingComment() tokenNodeOption {
	return func(o *tokenNodeOpt) {
		o.ignoreLeadingComment = true
	}
}

func ignoreComment() tokenNodeOption {
	return func(o *tokenNodeOpt) {
		o.ignoreHeadComment = true
		o.ignoreLeadingComment = true
	}
}

func withTokenNodePrefix(prefix ...string) tokenNodeOption {
	return func(o *tokenNodeOpt) {
		for _, p := range prefix {
			o.prefix = p
		}
	}

}
func withTokenNodeInfix(infix string) tokenNodeOption {
	return func(o *tokenNodeOpt) {
		o.infix = infix
	}
}

func expectSameLine() Option {
	return func(o *option) {
		o.mode = ModeExpectInSameLine
	}
}

func expectIndentInfix() Option {
	return func(o *option) {
		o.infix = Indent
	}
}

func withNode(nodes ...Node) Option {
	return func(o *option) {
		o.nodes = nodes
	}
}

func withMode(mode WriteMode) Option {
	return func(o *option) {
		o.mode = mode
	}
}

func withPrefix(prefix ...string) Option {
	return func(o *option) {
		for _, p := range prefix {
			o.prefix = p
		}
	}
}

func withInfix(infix string) Option {
	return func(o *option) {
		o.infix = infix
	}
}

func withRawText() Option {
	return func(o *option) {
		o.rawText = true
	}
}

// NewWriter returns a new Writer.
func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		tw:     tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		writer: writer,
	}
}

// NewBufferWriter returns a new buffer Writer.
func NewBufferWriter() *Writer {
	writer := bytes.NewBuffer(nil)
	return &Writer{
		tw:     tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		writer: writer,
	}
}

// String returns the string of the buffer.
func (w *Writer) String() string {
	buffer, ok := w.writer.(*bytes.Buffer)
	if !ok {
		return ""
	}
	w.Flush()
	return buffer.String()
}

// Flush flushes the buffer.
func (w *Writer) Flush() {
	_ = w.tw.Flush()
}

// NewLine writes a new line.
func (w *Writer) NewLine() {
	_, _ = fmt.Fprint(w.tw, NewLine)
}

// Write writes the node.
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

// WriteText writes the text.
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
		mode := opt.mode
		preIdx := idx - 1
		var preNodeHasLeading bool
		if preIdx > -1 && preIdx < len(opt.nodes) {
			preNode := opt.nodes[preIdx]
			preNodeHasLeading = preNode.HasLeadingCommentGroup()
		}
		if node.HasHeadCommentGroup() || preNodeHasLeading {
			mode = ModeAuto
		}

		if mode == ModeAuto && node.Pos().Line > line {
			textList = append(textList, NewLine)
		}
		line = node.End().Line
		if util.TrimWhiteSpace(node.Format()) == "" {
			continue
		}

		textList = append(textList, node.Format(opt.prefix))
	}

	text := strings.Join(textList, opt.infix)
	text = strings.ReplaceAll(text, " \n", "\n")
	text = strings.ReplaceAll(text, "\n ", "\n")
	if opt.rawText {
		_, _ = fmt.Fprint(w.writer, text)
		return
	}
	_, _ = fmt.Fprint(w.tw, text)
}
