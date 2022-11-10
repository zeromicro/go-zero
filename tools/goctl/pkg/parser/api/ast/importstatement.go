package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type ImportStmt interface {
	Stmt
	importNode()
}

type ImportLiteralStmt struct {
	Import token.Token
	Value  token.Token
}

func (i *ImportLiteralStmt) importNode() {}

func (i *ImportLiteralStmt) Pos() token.Position {
	return i.Import.Position
}

func (i *ImportLiteralStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, i.Import, i.Value)
	return w.String()
}

func (i *ImportLiteralStmt) stmtNode() {}

type ImportGroupStmt struct {
	Import token.Token
	LParen token.Token
	Values []token.Token
	RParen token.Token
}

func (i *ImportGroupStmt) importNode() {}

func (i *ImportGroupStmt) Pos() token.Position {
	return i.Import.Position
}

func (i *ImportGroupStmt) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, i.Import, i.LParen)
	if len(i.Values) > 0 {
		w.NewLine()
	}
	for _, v := range i.Values {
		w.Writeln(prefix, indent, v)
	}
	if len(i.Values) > 0 {
		w.Write(prefix, i.RParen)
	} else {
		w.Write(i.RParen)
	}

	return w.String()
}

func (i *ImportGroupStmt) stmtNode() {}
