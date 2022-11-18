package ast

import (
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

type ImportStmt interface {
	Stmt
	importNode()
}

type ImportLiteralStmt struct {
	Import *TokenNode
	Value  *TokenNode
}

func (i *ImportLiteralStmt) Format(prefix ...string) (result string) {
	if i.Value.IsZeroString() {
		return ""
	}
	w := NewBufferWriter()
	importNode := transferTokenNode(i.Import, ignoreLeadingComment(), withTokenNodePrefix(prefix...))
	w.Write(WithNode(importNode, i.Value), WithMode(ModeExpectInSameLine))
	return w.String()
}

func (i *ImportLiteralStmt) End() token.Position {
	return i.Value.End()
}

func (i *ImportLiteralStmt) importNode() {}

func (i *ImportLiteralStmt) Pos() token.Position {
	return i.Import.Pos()
}

func (i *ImportLiteralStmt) stmtNode() {}

type ImportGroupStmt struct {
	Import *TokenNode
	LParen *TokenNode
	Values []*TokenNode
	RParen *TokenNode
}

func (i *ImportGroupStmt) Format(prefix ...string) string {
	var textList []string
	for _, v := range i.Values {
		if v.IsZeroString() {
			continue
		}
		textList = append(textList, v.Format(Indent))
	}
	if len(textList) == 0 {
		return ""
	}

	importNode := transferTokenNode(i.Import, ignoreLeadingComment(), withTokenNodePrefix(prefix...))
	w := NewBufferWriter()
	w.Write(WithNode(importNode, i.LParen), expectSameLine())
	w.NewLine()
	for _, v := range i.Values {
		node := transferTokenNode(v, withTokenNodePrefix(peekOne(prefix)+Indent))
		w.Write(WithNode(node), expectSameLine())
		w.NewLine()
	}
	w.Write(WithNode(transferTokenNode(i.RParen, withTokenNodePrefix(prefix...))))
	return w.String()
}

func (i *ImportGroupStmt) End() token.Position {
	return i.RParen.End()
}

func (i *ImportGroupStmt) importNode() {}

func (i *ImportGroupStmt) Pos() token.Position {
	return i.Import.Pos()
}

func (i *ImportGroupStmt) stmtNode() {}
