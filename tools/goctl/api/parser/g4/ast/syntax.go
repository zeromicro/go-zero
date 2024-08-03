package ast

import "github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"

// SyntaxExpr describes syntax for api
type SyntaxExpr struct {
	Syntax      Expr
	Assign      Expr
	Version     Expr
	DocExpr     []Expr
	CommentExpr Expr
}

// VisitSyntaxLit implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitSyntaxLit(ctx *api.SyntaxLitContext) any {
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

// Format provides a formatter for api command, now nothing to do
func (s *SyntaxExpr) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two SyntaxExpr are equal
func (s *SyntaxExpr) Equal(v any) bool {
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

// Doc returns the document of SyntaxExpr, like // some text
func (s *SyntaxExpr) Doc() []Expr {
	return s.DocExpr
}

// Comment returns the comment of SyntaxExpr, like // some text
func (s *SyntaxExpr) Comment() Expr {
	return s.CommentExpr
}
