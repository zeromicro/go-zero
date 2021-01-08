package ast

import (
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
)

type SyntaxExpr struct {
	Syntax      Expr
	Assign      Expr
	Version     Expr
	DocExpr     []Expr
	CommentExpr Expr
}

func (v *ApiVisitor) VisitSyntaxLit(ctx *api.SyntaxLitContext) interface{} {
	syntax := v.newExprWithToken(ctx.GetSyntaxToken())
	assign := v.newExprWithToken(ctx.GetAssign())
	version := v.newExprWithToken(ctx.GetVersion())
	return &SyntaxExpr{
		Syntax:      syntax,
		Assign:      assign,
		Version:     version,
		DocExpr:     v.getDoc(ctx),
		CommentExpr: v.getComment(ctx),
	}
}

func (s *SyntaxExpr) Format() error {
	// todo
	return nil
}

func (s *SyntaxExpr) Equal(v interface{}) bool {
	if v == nil {
		return false
	}

	syntax, ok := v.(*SyntaxExpr)
	if !ok {
		return false
	}

	if !EqualDoc(s, syntax) {
		return false
	}

	return s.Syntax.Equal(syntax.Syntax) &&
		s.Assign.Equal(syntax.Assign) &&
		s.Version.Equal(syntax.Version)
}

func (s *SyntaxExpr) Doc() []Expr {
	return s.DocExpr
}

func (s *SyntaxExpr) Comment() Expr {
	return s.CommentExpr
}
