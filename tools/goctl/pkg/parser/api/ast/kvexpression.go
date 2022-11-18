package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type KVExpr struct {
	Key   *TokenNode
	Value *TokenNode
}

func (i *KVExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.Write(WithNode(i.Key, i.Value), WithPrefix(prefix...), WithInfix(Indent))
	return w.String()
}

func (i *KVExpr) End() token.Position {
	return i.Value.End()
}

func (i *KVExpr) Pos() token.Position {
	return i.Key.Pos()
}

func (i *KVExpr) exprNode() {}
