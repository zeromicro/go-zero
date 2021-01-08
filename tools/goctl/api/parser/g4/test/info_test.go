package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
)

var infoAccept = func(p *api.ApiParserParser, visitor *ast.ApiVisitor) interface{} {
	return p.InfoSpec().Accept(visitor)
}

func TestInfo(t *testing.T) {
	t.Run("matched", func(t *testing.T) {
		v, err := parser.Accept(infoAccept, `
			info(
				title: foo
			)
		`)

		assert.Nil(t, err)
		info := v.(*ast.InfoExpr)
		assert.True(t, info.Equal(&ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Kvs: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("title"),
					Value: ast.NewTextExpr("foo"),
				},
			},
		}))

		v, err = parser.Accept(infoAccept, `
			info(
				title: 中文(bar)
			)
		`)
		assert.Nil(t, err)
		info = v.(*ast.InfoExpr)
		assert.True(t, info.Equal(&ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Kvs: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("title"),
					Value: ast.NewTextExpr("中文(bar)"),
				},
			},
		}))

		v, err = parser.Accept(infoAccept, `
			info(
				foo: "new
line"
			)
		`)
		assert.Nil(t, err)
		info = v.(*ast.InfoExpr)
		assert.True(t, info.Equal(&ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Kvs: []*ast.KvExpr{
				{
					Key: ast.NewTextExpr("foo"),
					Value: ast.NewTextExpr(`"new
line"`),
				},
			},
		}))
	})

	t.Run("matched doc", func(t *testing.T) {
		v, err := parser.Accept(infoAccept, `
			// doc
			info( // comment
				// foo
				title: foo // bar
			)
		`)
		assert.Nil(t, err)
		info := v.(*ast.InfoExpr)
		assert.True(t, info.Equal(&ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Kvs: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("title"),
					Value: ast.NewTextExpr("foo"),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// foo"),
					},
					CommentExpr: ast.NewTextExpr("// bar"),
				},
			},
		}))

		v, err = parser.Accept(infoAccept, `
			/**doc block*/
			info( /**line block*/
				/**foo*/
				title: foo /*bar**/
			)
		`)
		assert.Nil(t, err)
		info = v.(*ast.InfoExpr)
		assert.True(t, info.Equal(&ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Kvs: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("title"),
					Value: ast.NewTextExpr("foo"),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("/**foo*/"),
					},
					CommentExpr: ast.NewTextExpr("/*bar**/"),
				},
			},
		}))

		v, err = parser.Accept(infoAccept, `
			info( 
				// doc
				title: foo 
				// doc
				author: bar
			)
		`)
		assert.Nil(t, err)
		info = v.(*ast.InfoExpr)
		assert.True(t, info.Equal(&ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Kvs: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("title"),
					Value: ast.NewTextExpr("foo"),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// doc"),
					},
				},
				{
					Key:   ast.NewTextExpr("author"),
					Value: ast.NewTextExpr("bar"),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// doc"),
					},
				},
			},
		}))

	})

	t.Run("mismatched", func(t *testing.T) {
		_, err := parser.Accept(infoAccept, `
			info(
				title
			)
		`)
		assert.Error(t, err)

		_, err = parser.Accept(infoAccept, `
			info(
				:title
			)
		`)
		assert.Error(t, err)

		_, err = parser.Accept(infoAccept, `
			info(
				foo bar
			)
		`)
		assert.Error(t, err)

		_, err = parser.Accept(infoAccept, `
			info(
				foo : new 
line
			)
		`)
		assert.Error(t, err)

		_, err = parser.Accept(infoAccept, `
			info()
		`)
		assert.Error(t, err)
	})
}
