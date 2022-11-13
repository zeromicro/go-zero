package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type InfoStmt struct {
	Info   token.Token
	LParen token.Token
	Values []*KVExpr
	RParen token.Token
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Position
}

func (i *InfoStmt) stmtNode() {}
