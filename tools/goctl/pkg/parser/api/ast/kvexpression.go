package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type KVExpr struct {
	Key   token.Token
	Value token.Token
}

func (i *KVExpr) Pos() token.Position {
	return i.Key.Position
}

func (i *KVExpr) exprNode() {}
