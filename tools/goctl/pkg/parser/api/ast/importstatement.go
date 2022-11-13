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

func (i *ImportGroupStmt) stmtNode() {}
