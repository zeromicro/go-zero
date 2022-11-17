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

func (i *ImportLiteralStmt) Format(...string) (result string) {
	if i.Value.IsZeroString() {
		return ""
	}
	w := NewBufferWriter()
	w.Write(WithNode(i.Import, i.Value), WithMode(ModeExpectInSameLine))
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

func (i *ImportGroupStmt) Format(...string) string {
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

	w := NewBufferWriter()
	w.Write(WithNode(i.Import, i.LParen), WithMode(ModeExpectInSameLine))
	w.NewLine()
	for _, v := range i.Values {
		if v.IsZeroString() {
			continue
		}
		w.Write(WithNode(v), WithPrefix(Indent))
		w.NewLine()
	}
	w.WriteText(i.RParen.Format())
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
