package ast

import "github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"

// InfoExpr defines info syntax for api
type InfoExpr struct {
	Info Expr
	Lp   Expr
	Rp   Expr
	Kvs  []*KvExpr
}

// VisitInfoSpec implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitInfoSpec(ctx *api.InfoSpecContext) any {
	var expr InfoExpr
	expr.Info = v.newExprWithToken(ctx.GetInfoToken())
	expr.Lp = v.newExprWithToken(ctx.GetLp())
	expr.Rp = v.newExprWithToken(ctx.GetRp())
	list := ctx.AllKvLit()
	for _, each := range list {
		kvExpr := each.Accept(v).(*KvExpr)
		expr.Kvs = append(expr.Kvs, kvExpr)
	}

	if v.infoFlag {
		v.panic(expr.Info, "duplicate declaration 'info'")
	}

	return &expr
}

// Format provides a formatter for api command, now nothing to do
func (i *InfoExpr) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two InfoExpr are equal
func (i *InfoExpr) Equal(v any) bool {
	if v == nil {
		return false
	}

	info, ok := v.(*InfoExpr)
	if !ok {
		return false
	}

	if !i.Info.Equal(info.Info) {
		return false
	}

	var expected, actual []*KvExpr
	expected = append(expected, i.Kvs...)
	actual = append(actual, info.Kvs...)

	if len(expected) != len(actual) {
		return false
	}

	for index, each := range expected {
		ac := actual[index]
		if !each.Equal(ac) {
			return false
		}
	}

	return true
}
