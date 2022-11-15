package ast

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/placeholder"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

const NilIndent = ""
const WhiteSpace = " "
const Indent = "\t"
const NewLine = "\n"

type Writer struct {
	tw            *tabwriter.Writer
	lastWriteNode Node
	nodeSet       *NodeSet

	writer io.Writer
	skip   map[Node]placeholder.Type
}

func NewWriter(writer io.Writer, tokenSet *NodeSet) *Writer {
	return &Writer{
		tw:            tabwriter.NewWriter(writer, 1, 8, 1, ' ', tabwriter.TabIndent),
		lastWriteNode: initNode,
		nodeSet:       tokenSet,
		writer:        writer,
		skip:          make(map[Node]placeholder.Type),
	}
}

func (w *Writer) Fork() *Writer {
	var buffer = bytes.NewBuffer(nil)
	return NewWriter(buffer, w.nodeSet)
}

func (w *Writer) String() string {
	w.Flush()
	if bw, ok := w.writer.(*bytes.Buffer); ok {
		return bw.String()
	}
	return ""
}

func (w *Writer) WriteBetween(prefix string, left, right Node) {
	nodes := w.nodeSet.Between(left, right, AllIn)
	if len(nodes) > 0 {
		w.Write(prefix, nodes...)
	}
}

func (w *Writer) WriteSpaceInfixBetween(prefix string, left, right Node) {
	nodes := w.nodeSet.Between(left, right, AllIn)
	if len(nodes) > 0 {
		w.WriteSpaceInfix(prefix, nodes...)
	}
}

func (w *Writer) WriteSpaceInfix(prefix string, nodes ...Node) {
	nodes = w.filterSkipNode(nodes)
	if len(nodes) == 0 {
		return
	}

	defer func() {
		tail := nodes[len(nodes)-1]
		lineAfter := w.nodeSet.LineCommentAfter(tail)
		if len(lineAfter) > 0 {
			w.write(NilIndent, lineAfter...)
		}
	}()
	var hasDoc = false
	for _, e := range nodes {
		if isComment(e) {
			hasDoc = true
			break
		}
	}

	var one = nodes[0]
	gaps := w.nodeSet.Between(w.lastWriteNode, one, NotIn)
	if len(gaps) > 0 {
		w.write(NilIndent, gaps...)
	}

	_, _ = fmt.Fprint(w.tw, prefix)
	for idx, node := range nodes {
		if node.Pos().Line > w.lastWriteNode.End().Line && (idx == 0 || hasDoc) {
			w.NewLine()
		}
		_, _ = fmt.Fprint(w.tw, node.Format())
		if idx < len(nodes)-1 {
			_, _ = fmt.Fprint(w.tw, WhiteSpace)
		}
		w.lastWriteNode = node
	}
}

func (w *Writer) WriteInOneLine(prefix string, nodes ...Node) {
	nodes = w.filterSkipNode(nodes)
	if len(nodes) == 0 {
		return
	}
	defer func() {
		tail := nodes[len(nodes)-1]
		lineAfter := w.nodeSet.LineCommentAfter(tail)
		if len(lineAfter) > 0 {
			w.write(NilIndent, lineAfter...)
		}
	}()
	var one = nodes[0]
	gaps := w.nodeSet.Between(w.lastWriteNode, one, NotIn)
	if len(gaps) > 0 {
		w.write(NilIndent, gaps...)
	}

	var lastWriteToken = w.lastWriteNode
	var hasDocument = false
	var list []string
	for _, node := range nodes {
		if isComment(node) {
			hasDocument = true
		}
		list = append(list, node.Format())
		lastWriteToken = node
	}
	if !hasDocument {
		_, _ = fmt.Fprint(w.tw, prefix)
		_, _ = fmt.Fprint(w.tw, strings.Join(list, WhiteSpace))
		w.lastWriteNode = lastWriteToken
		return
	}

	w.write(prefix, nodes...)
}

func (w *Writer) Write(prefix string, nodes ...Node) {
	nodes = w.filterSkipNode(nodes)
	if len(nodes) == 0 {
		return
	}

	defer func() {
		tail := nodes[len(nodes)-1]
		lineAfter := w.nodeSet.LineCommentAfter(tail)
		if len(lineAfter) > 0 {
			w.write(NilIndent, lineAfter...)
		}
	}()
	var one = nodes[0]
	gaps := w.nodeSet.Between(w.lastWriteNode, one, NotIn)
	if len(gaps) > 0 {
		w.write(NilIndent, gaps...)
	}

	var lastWriteToken = w.lastWriteNode
	var inOneLine = true
	var list []string
	for _, node := range nodes {
		if one.Pos().Line != node.End().Line {
			inOneLine = false
		}
		list = append(list, node.Format())
		lastWriteToken = node
	}
	if inOneLine && len(list) > 0 {
		if one.Pos().Line > w.lastWriteNode.End().Line {
			w.NewLine()
		}
		_, _ = fmt.Fprint(w.tw, prefix)
		_, _ = fmt.Fprint(w.tw, strings.Join(list, Indent))
		w.lastWriteNode = lastWriteToken
		return
	}

	w.write(prefix, nodes...)
}

func (w *Writer) filterSkipNode(nodes []Node) []Node {
	var list []Node
	for _, node := range nodes {
		if w.canSkip(node) {
			continue
		}
		list = append(list, node)
	}
	return list
}

func (w *Writer) canSkip(node Node) bool {
	tokenNode, ok := node.(*TokenNode)
	if ok && tokenNode.Token.IsType(token.EOF) {
		return true
	}
	_, ok = w.skip[node]
	return ok
}

func (w *Writer) write(prefix string, nodes ...Node) {
	nodes = w.filterSkipNode(nodes)
	if len(nodes) == 0 {
		return
	}
	_, _ = fmt.Fprint(w.tw, prefix)
	for idx, node := range nodes {
		if node.Pos().Line > w.lastWriteNode.End().Line {
			w.NewLine()
		}
		_, _ = fmt.Fprint(w.tw, node.Format())
		if idx < len(nodes)-1 {
			_, _ = fmt.Fprint(w.tw, WhiteSpace)
		}
		w.lastWriteNode = node
	}
}

func (w *Writer) Skip(nodes ...Node) {
	for _, node := range nodes {
		w.skip[node] = placeholder.PlaceHolder
	}
}

func (w *Writer) NewLine() {
	_, _ = fmt.Fprint(w.tw, NewLine)
}

func (w *Writer) WriteTailGaps() {
	list := w.nodeSet.CommentBetween(w.lastWriteNode, w.nodeSet.LastToken(), RightIn)
	w.Write(NilIndent, list...)
	w.NewLine()
}

func (w *Writer) Flush() {
	_ = w.tw.Flush()
}

func isComment(node Node) bool {
	_, ok := node.(*CommentStmt)
	return ok
}
