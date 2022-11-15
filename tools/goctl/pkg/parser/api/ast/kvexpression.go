package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type KVExpr struct {
	Key   token.Token
	Value token.Token

	fw *Writer
}

func (i *KVExpr) Format(prefix ...string) string {
	//TODO implement me
	panic("implement me")
}

func (i *KVExpr) End() token.Position {
	return i.Value.Position
}

func (i *KVExpr) Pos() token.Position {
	return i.Key.Position
}

func (i *KVExpr) exprNode() {}
