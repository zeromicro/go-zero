package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type ImportStmt interface {
	Stmt
	importNode()
}

type ImportLiteralStmt struct {
	Import token.Token
	Value  token.Token

	fw *Writer
}

func (i *ImportLiteralStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (i *ImportLiteralStmt) End() token.Position {
	return i.Value.Position
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

	nodeSet *NodeSet
}

func (i *ImportGroupStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (i *ImportGroupStmt) End() token.Position {
	return i.RParen.Position
}

func (i *ImportGroupStmt) importNode() {}

func (i *ImportGroupStmt) Pos() token.Position {
	return i.Import.Position
}

func (i *ImportGroupStmt) stmtNode() {}
