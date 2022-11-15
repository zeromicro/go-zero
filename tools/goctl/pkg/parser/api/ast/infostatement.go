package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type InfoStmt struct {
	Info   token.Token
	LParen token.Token
	Values []*KVExpr
	RParen token.Token

	fw *Writer
}

func (i *InfoStmt) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (i *InfoStmt) End() token.Position {
	return i.RParen.Position
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Position
}

func (i *InfoStmt) stmtNode() {}
