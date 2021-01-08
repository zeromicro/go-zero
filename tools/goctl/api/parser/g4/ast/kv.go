package ast

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
)

type KvExpr struct {
	Key         Expr
	Value       Expr
	DocExpr     []Expr
	CommentExpr Expr
}

func (v *ApiVisitor) VisitKvLit(ctx *api.KvLitContext) interface{} {
	var kvExpr KvExpr
	kvExpr.Key = v.newExprWithToken(ctx.GetKey())
	commentExpr := v.getComment(ctx)
	if ctx.GetValue() != nil {
		valueText := ctx.GetValue().GetText()
		valueExpr := v.newExprWithToken(ctx.GetValue())
		if strings.Contains(valueText, "//") {
			if commentExpr == nil {
				commentExpr = v.newExprWithToken(ctx.GetValue())
				commentExpr.SetText("")
			}

			index := strings.Index(valueText, "//")
			commentExpr.SetText(valueText[index:])
			valueExpr.SetText(strings.TrimSpace(valueText[:index]))
		} else if strings.Contains(valueText, "/*") {
			if commentExpr == nil {
				commentExpr = v.newExprWithToken(ctx.GetValue())
				commentExpr.SetText("")
			}

			index := strings.Index(valueText, "/*")
			commentExpr.SetText(valueText[index:])
			valueExpr.SetText(strings.TrimSpace(valueText[:index]))
		}

		kvExpr.Value = valueExpr
	}

	kvExpr.DocExpr = v.getDoc(ctx)
	kvExpr.CommentExpr = commentExpr
	return &kvExpr
}

func (k *KvExpr) Format() error {
	// todo
	return nil
}

func (k *KvExpr) Equal(v interface{}) bool {
	if v == nil {
		return false
	}

	kv, ok := v.(*KvExpr)
	if !ok {
		return false
	}

	if !EqualDoc(k, kv) {
		return false
	}

	return k.Key.Equal(kv.Key) && k.Value.Equal(kv.Value)
}

func (k *KvExpr) Doc() []Expr {
	return k.DocExpr
}

func (k *KvExpr) Comment() Expr {
	return k.CommentExpr
}
