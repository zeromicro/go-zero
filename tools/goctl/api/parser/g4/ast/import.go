package ast

import "github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"

// ImportExpr defines import syntax for api
type ImportExpr struct {
	Import      Expr
	Value       Expr
	DocExpr     []Expr
	CommentExpr Expr
}

// VisitImportSpec implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitImportSpec(ctx *api.ImportSpecContext) any {
	var list []*ImportExpr
	if ctx.ImportLit() != nil {
		lits := ctx.ImportLit().Accept(v).([]*ImportExpr)
		list = append(list, lits...)
	}
	if ctx.ImportBlock() != nil {
		blocks := ctx.ImportBlock().Accept(v).([]*ImportExpr)
		list = append(list, blocks...)
	}

	return list
}

// VisitImportLit implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitImportLit(ctx *api.ImportLitContext) any {
	importToken := v.newExprWithToken(ctx.GetImportToken())
	valueExpr := ctx.ImportValue().Accept(v).(Expr)
	return []*ImportExpr{
		{
			Import:      importToken,
			Value:       valueExpr,
			DocExpr:     v.getDoc(ctx),
			CommentExpr: v.getComment(ctx),
		},
	}
}

// VisitImportBlock implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitImportBlock(ctx *api.ImportBlockContext) any {
	importToken := v.newExprWithToken(ctx.GetImportToken())
	values := ctx.AllImportBlockValue()
	var list []*ImportExpr

	for _, value := range values {
		importExpr := value.Accept(v).(*ImportExpr)
		importExpr.Import = importToken
		list = append(list, importExpr)
	}

	return list
}

// VisitImportBlockValue implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitImportBlockValue(ctx *api.ImportBlockValueContext) any {
	value := ctx.ImportValue().Accept(v).(Expr)
	return &ImportExpr{
		Value:       value,
		DocExpr:     v.getDoc(ctx),
		CommentExpr: v.getComment(ctx),
	}
}

// VisitImportValue implements from api.BaseApiParserVisitor
func (v *ApiVisitor) VisitImportValue(ctx *api.ImportValueContext) any {
	return v.newExprWithTerminalNode(ctx.STRING())
}

// Format provides a formatter for api command, now nothing to do
func (i *ImportExpr) Format() error {
	// todo
	return nil
}

// Equal compares whether the element literals in two ImportExpr are equal
func (i *ImportExpr) Equal(v any) bool {
	if v == nil {
		return false
	}

	imp, ok := v.(*ImportExpr)
	if !ok {
		return false
	}

	if !EqualDoc(i, imp) {
		return false
	}

	return i.Import.Equal(imp.Import) && i.Value.Equal(imp.Value)
}

// Doc returns the document of ImportExpr, like // some text
func (i *ImportExpr) Doc() []Expr {
	return i.DocExpr
}

// Comment returns the comment of ImportExpr, like // some text
func (i *ImportExpr) Comment() Expr {
	return i.CommentExpr
}
