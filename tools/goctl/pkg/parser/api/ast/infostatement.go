package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type InfoStmt struct {
	Info   *TokenNode
	LParen *TokenNode
	Values []*KVExpr
	RParen *TokenNode
}

func (i *InfoStmt) Format(prefix ...string) (result string) {
	return
}

func (i *InfoStmt) End() token.Position {
	return i.RParen.End()
}

func (i *InfoStmt) Pos() token.Position {
	return i.Info.Pos()
}

func (i *InfoStmt) stmtNode() {}
