package ast

import "github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"

type KVExpr struct {
	Key   token.Token
	Colon token.Token
	Value token.Token
}

func (i *KVExpr) Pos() token.Position {
	return i.Key.Position
}

func (i *KVExpr) Format(prefix string) string {
	w := NewWriter()
	w.WriteWithWhiteSpaceInfix(prefix, i.Key.Text+i.Colon.Text, i.Value)
	return w.String()
}

func (i *KVExpr) exprNode() {}
