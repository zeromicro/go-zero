package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
)

var syntaxAccept = func(p *api.ApiParserParser, visitor *ast.ApiVisitor) interface{} {
	return p.SyntaxLit().Accept(visitor)
}

func TestSyntax(t *testing.T) {
	t.Run("matched", func(t *testing.T) {
		v, err := parser.Accept(syntaxAccept, `syntax = "v1"`)
		assert.Nil(t, err)

		syntax := v.(*ast.SyntaxExpr)
		assert.True(t, syntax.Equal(&ast.SyntaxExpr{
			Syntax:  ast.NewTextExpr("syntax"),
			Assign:  ast.NewTextExpr("="),
			Version: ast.NewTextExpr(`"v1"`),
		}))
	})

	t.Run("expecting syntax", func(t *testing.T) {
		_, err := parser.Accept(syntaxAccept, `= "v1"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syn = "v1"`)
		assert.Error(t, err)
	})

	t.Run("missing assign", func(t *testing.T) {
		_, err := parser.Accept(syntaxAccept, `syntax  "v1"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax + "v1"`)
		assert.Error(t, err)
	})

	t.Run("mismatched version", func(t *testing.T) {
		_, err := parser.Accept(syntaxAccept, `syntax="v0"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax = "v1a"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax = "vv1"`)
		assert.Error(t, err)

		_, err = parser.Accept(syntaxAccept, `syntax = "1"`)
		assert.Error(t, err)
	})

	t.Run("with comment", func(t *testing.T) {
		v, err := parser.Accept(syntaxAccept, `
		// doc
		syntax="v1" // line comment`)
		assert.Nil(t, err)

		syntax := v.(*ast.SyntaxExpr)
		assert.True(t, syntax.Equal(&ast.SyntaxExpr{
			Syntax:  ast.NewTextExpr("syntax"),
			Assign:  ast.NewTextExpr("="),
			Version: ast.NewTextExpr(`"v1"`),
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// doc"),
			},
			CommentExpr: ast.NewTextExpr("// line comment"),
		}))
	})
}
