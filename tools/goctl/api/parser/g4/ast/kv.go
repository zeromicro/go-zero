package ast

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
)

// KvExpr describes key-value for api
type KvExpr struct {
	Key         Expr
	Value       Expr
	DocExpr     []Expr
	CommentExpr Expr
}

// VisitKvLit implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitKvLit(ctx *api.KvLitContext) any {
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

// Format provides a formatter for api command, now nothing to do
func (k *KvExpr) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two KvExpr are equal
func (k *KvExpr) Equal(v any) bool {
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

// Doc returns the document of KvExpr, like // some text
func (k *KvExpr) Doc() []Expr {
	return k.DocExpr
}

// Comment returns the comment of KvExpr, like // some text
func (k *KvExpr) Comment() Expr {
	return k.CommentExpr
}
