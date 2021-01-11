package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/gen/api"
)

func TestBody(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.Body().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `(Foo)`)
		assert.Nil(t, err)
		body := v.(*ast.Body)
		assert.True(t, body.Equal(&ast.Body{
			Lp:   ast.NewTextExpr("("),
			Rp:   ast.NewTextExpr(")"),
			Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `(var)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `()`)
		assert.Nil(t, err)
	})
}

func TestRoute(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.Route().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `post /foo/foo-bar/:bar (Foo) returns (Bar)`)
		assert.Nil(t, err)
		route := v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method: ast.NewTextExpr("post"),
			Path:   ast.NewTextExpr("/foo/foo-bar/:bar"),
			Req: &ast.Body{
				Lp:   ast.NewTextExpr("("),
				Rp:   ast.NewTextExpr(")"),
				Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
			},
			ReturnToken: ast.NewTextExpr("returns"),
			Reply: &ast.Body{
				Lp:   ast.NewTextExpr("("),
				Rp:   ast.NewTextExpr(")"),
				Name: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar (Foo)`)
		assert.Nil(t, err)
		route = v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method: ast.NewTextExpr("post"),
			Path:   ast.NewTextExpr("/foo/foo-bar/:bar"),
			Req: &ast.Body{
				Lp:   ast.NewTextExpr("("),
				Rp:   ast.NewTextExpr(")"),
				Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar returns (Bar)`)
		assert.Nil(t, err)
		route = v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method:      ast.NewTextExpr("post"),
			Path:        ast.NewTextExpr("/foo/foo-bar/:bar"),
			ReturnToken: ast.NewTextExpr("returns"),
			Reply: &ast.Body{
				Lp:   ast.NewTextExpr("("),
				Rp:   ast.NewTextExpr(")"),
				Name: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar returns ([]Bar)`)
		assert.Nil(t, err)
		route = v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method:      ast.NewTextExpr("post"),
			Path:        ast.NewTextExpr("/foo/foo-bar/:bar"),
			ReturnToken: ast.NewTextExpr("returns"),
			Reply: &ast.Body{
				Lp: ast.NewTextExpr("("),
				Rp: ast.NewTextExpr(")"),
				Name: &ast.Array{
					ArrayExpr: ast.NewTextExpr("[]Bar"),
					LBrack:    ast.NewTextExpr("["),
					RBrack:    ast.NewTextExpr("]"),
					Literal:   &ast.Literal{Literal: ast.NewTextExpr("Bar")},
				},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar returns ([]*Bar)`)
		assert.Nil(t, err)
		route = v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method:      ast.NewTextExpr("post"),
			Path:        ast.NewTextExpr("/foo/foo-bar/:bar"),
			ReturnToken: ast.NewTextExpr("returns"),
			Reply: &ast.Body{
				Lp: ast.NewTextExpr("("),
				Rp: ast.NewTextExpr(")"),
				Name: &ast.Array{
					ArrayExpr: ast.NewTextExpr("[]*Bar"),
					LBrack:    ast.NewTextExpr("["),
					RBrack:    ast.NewTextExpr("]"),
					Literal: &ast.Pointer{
						PointerExpr: ast.NewTextExpr("*Bar"),
						Star:        ast.NewTextExpr("*"),
						Name:        ast.NewTextExpr("Bar"),
					}},
			},
		}))

		v, err = parser.Accept(fn, `post /foo/foo-bar/:bar`)
		assert.Nil(t, err)
		route = v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method: ast.NewTextExpr("post"),
			Path:   ast.NewTextExpr("/foo/foo-bar/:bar"),
		}))

		v, err = parser.Accept(fn, `
		// foo
		post /foo/foo-bar/:bar // bar`)
		assert.Nil(t, err)
		route = v.(*ast.Route)
		assert.True(t, route.Equal(&ast.Route{
			Method: ast.NewTextExpr("post"),
			Path:   ast.NewTextExpr("/foo/foo-bar/:bar"),
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// foo"),
			},
			CommentExpr: ast.NewTextExpr("// bar"),
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `posts /foo`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `gets /foo`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post /foo/:`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post /foo/`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post foo/bar`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `post /foo/bar return (Bar)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` /foo/bar returns (Bar)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post   returns (Bar)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns (int)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns (*int)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns ([]var)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, ` post /foo/bar returns (const)`)
		assert.Error(t, err)
	})
}

func TestAtHandler(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.AtHandler().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `@handler foo`)
		assert.Nil(t, err)
		atHandler := v.(*ast.AtHandler)
		assert.True(t, atHandler.Equal(&ast.AtHandler{
			AtHandlerToken: ast.NewTextExpr("@handler"),
			Name:           ast.NewTextExpr("foo"),
		}))

		v, err = parser.Accept(fn, `
		// foo
		@handler foo // bar`)
		assert.Nil(t, err)
		atHandler = v.(*ast.AtHandler)
		assert.True(t, atHandler.Equal(&ast.AtHandler{
			AtHandlerToken: ast.NewTextExpr("@handler"),
			Name:           ast.NewTextExpr("foo"),
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// foo"),
			},
			CommentExpr: ast.NewTextExpr("// bar"),
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, ``)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@handler`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@handler "foo"`)
		assert.Error(t, err)
	})

}

func TestAtDoc(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.AtDoc().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `@doc "foo"`)
		assert.Nil(t, err)
		atDoc := v.(*ast.AtDoc)
		assert.True(t, atDoc.Equal(&ast.AtDoc{
			AtDocToken: ast.NewTextExpr("@doc"),
			LineDoc:    ast.NewTextExpr(`"foo"`),
		}))

		v, err = parser.Accept(fn, `@doc("foo")`)
		assert.Nil(t, err)
		atDoc = v.(*ast.AtDoc)
		assert.True(t, atDoc.Equal(&ast.AtDoc{
			AtDocToken: ast.NewTextExpr("@doc"),
			Lp:         ast.NewTextExpr("("),
			Rp:         ast.NewTextExpr(")"),
			LineDoc:    ast.NewTextExpr(`"foo"`),
		}))

		v, err = parser.Accept(fn, `@doc(
			foo: bar
		)`)
		assert.Nil(t, err)
		atDoc = v.(*ast.AtDoc)
		assert.True(t, atDoc.Equal(&ast.AtDoc{
			AtDocToken: ast.NewTextExpr("@doc"),
			Lp:         ast.NewTextExpr("("),
			Rp:         ast.NewTextExpr(")"),
			Kv: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("foo"),
					Value: ast.NewTextExpr("bar"),
				},
			},
		}))

		v, err = parser.Accept(fn, `@doc(
			// foo
			foo: bar // bar
		)`)
		assert.Nil(t, err)
		atDoc = v.(*ast.AtDoc)
		assert.True(t, atDoc.Equal(&ast.AtDoc{
			AtDocToken: ast.NewTextExpr("@doc"),
			Lp:         ast.NewTextExpr("("),
			Rp:         ast.NewTextExpr(")"),
			Kv: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("foo"),
					Value: ast.NewTextExpr("bar"),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// foo"),
					},
					CommentExpr: ast.NewTextExpr("// bar"),
				},
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `@doc("foo"`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@doc "foo")`)
		assert.Error(t, err)
	})
}

func TestServiceRoute(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.ServiceRoute().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `
		@doc("foo")
		// foo/bar
		// foo
		@handler foo // bar
		// foo/bar
		// foo
		post /foo (Foo) returns (Bar) // bar
		`)
		assert.Nil(t, err)
		sr := v.(*ast.ServiceRoute)
		assert.True(t, sr.Equal(&ast.ServiceRoute{
			AtDoc: &ast.AtDoc{
				AtDocToken: ast.NewTextExpr("@doc"),
				Lp:         ast.NewTextExpr("("),
				Rp:         ast.NewTextExpr(")"),
				LineDoc:    ast.NewTextExpr(`"foo"`),
			},
			AtHandler: &ast.AtHandler{
				AtHandlerToken: ast.NewTextExpr("@handler"),
				Name:           ast.NewTextExpr("foo"),
				DocExpr: []ast.Expr{
					ast.NewTextExpr("// foo"),
				},
				CommentExpr: ast.NewTextExpr("// bar"),
			},
			Route: &ast.Route{
				Method: ast.NewTextExpr("post"),
				Path:   ast.NewTextExpr("/foo"),
				Req: &ast.Body{
					Lp:   ast.NewTextExpr("("),
					Rp:   ast.NewTextExpr(")"),
					Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
				},
				ReturnToken: ast.NewTextExpr("returns"),
				Reply: &ast.Body{
					Lp:   ast.NewTextExpr("("),
					Rp:   ast.NewTextExpr(")"),
					Name: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
				},
				DocExpr: []ast.Expr{
					ast.NewTextExpr("// foo"),
				},
				CommentExpr: ast.NewTextExpr("// bar"),
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `post /foo (Foo) returns (Bar) // bar`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@handler foo`)
		assert.Error(t, err)
	})
}

func TestServiceApi(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.ServiceApi().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `
		service foo-api{
			@doc("foo")
			// foo/bar
			// foo
			@handler foo // bar
			// foo/bar
			// foo
			post /foo (Foo) returns (Bar) // bar
		}
		`)
		assert.Nil(t, err)
		api := v.(*ast.ServiceApi)
		assert.True(t, api.Equal(&ast.ServiceApi{
			ServiceToken: ast.NewTextExpr("service"),
			Name:         ast.NewTextExpr("foo-api"),
			Lbrace:       ast.NewTextExpr("{"),
			Rbrace:       ast.NewTextExpr("}"),
			ServiceRoute: []*ast.ServiceRoute{
				{
					AtDoc: &ast.AtDoc{
						AtDocToken: ast.NewTextExpr("@doc"),
						Lp:         ast.NewTextExpr("("),
						Rp:         ast.NewTextExpr(")"),
						LineDoc:    ast.NewTextExpr(`"foo"`),
					},
					AtHandler: &ast.AtHandler{
						AtHandlerToken: ast.NewTextExpr("@handler"),
						Name:           ast.NewTextExpr("foo"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// foo"),
						},
						CommentExpr: ast.NewTextExpr("// bar"),
					},
					Route: &ast.Route{
						Method: ast.NewTextExpr("post"),
						Path:   ast.NewTextExpr("/foo"),
						Req: &ast.Body{
							Lp:   ast.NewTextExpr("("),
							Rp:   ast.NewTextExpr(")"),
							Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
						},
						ReturnToken: ast.NewTextExpr("returns"),
						Reply: &ast.Body{
							Lp:   ast.NewTextExpr("("),
							Rp:   ast.NewTextExpr(")"),
							Name: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
						},
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// foo"),
						},
						CommentExpr: ast.NewTextExpr("// bar"),
					},
				},
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `services foo-api{}`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `service foo-api{`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `service foo-api{
		post /foo
		}`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `service foo-api{
		@handler foo
		}`)
		assert.Error(t, err)
	})
}

func TestAtServer(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.AtServer().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, `
		@server(
			// foo
			foo1: bar1 // bar
			// foo
			foo2: "bar2" // bar
			/**foo*/
			foo3: "foo
			bar" /**bar*/		
		)
		`)
		assert.Nil(t, err)
		as := v.(*ast.AtServer)
		assert.True(t, as.Equal(&ast.AtServer{
			AtServerToken: ast.NewTextExpr("@server"),
			Lp:            ast.NewTextExpr("("),
			Rp:            ast.NewTextExpr(")"),
			Kv: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("foo1"),
					Value: ast.NewTextExpr("bar1"),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// foo"),
					},
					CommentExpr: ast.NewTextExpr("// bar"),
				},
				{
					Key:   ast.NewTextExpr("foo2"),
					Value: ast.NewTextExpr(`"bar2"`),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// foo"),
					},
					CommentExpr: ast.NewTextExpr("// bar"),
				},
				{
					Key: ast.NewTextExpr("foo3"),
					Value: ast.NewTextExpr(`"foo
			bar"`),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("/**foo*/"),
					},
					CommentExpr: ast.NewTextExpr("/**bar*/"),
				},
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `server (
			foo:bar
		)`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@server ()`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `@server (
			foo: bar
		`)
		assert.Error(t, err)
	})
}

func TestServiceSpec(t *testing.T) {
	fn := func(p *api.ApiParserParser, v *ast.ApiVisitor) interface{} {
		return p.ServiceSpec().Accept(v)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `
		service foo-api{
			@handler foo
			post /foo returns ([]int)
		}
		`)
		assert.Nil(t, err)

		v, err := parser.Accept(fn, `
		@server(
			// foo
			foo1: bar1 // bar
			// foo
			foo2: "bar2" // bar
			/**foo*/
			foo3: "foo
			bar" /**bar*/		
		)
		service foo-api{
			@doc("foo")
			// foo/bar
			// foo
			@handler foo // bar
			// foo/bar
			// foo
			post /foo (Foo) returns (Bar) // bar
		}
		`)
		assert.Nil(t, err)
		service := v.(*ast.Service)
		assert.True(t, service.Equal(&ast.Service{
			AtServer: &ast.AtServer{
				AtServerToken: ast.NewTextExpr("@server"),
				Lp:            ast.NewTextExpr("("),
				Rp:            ast.NewTextExpr(")"),
				Kv: []*ast.KvExpr{
					{
						Key:   ast.NewTextExpr("foo1"),
						Value: ast.NewTextExpr("bar1"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// foo"),
						},
						CommentExpr: ast.NewTextExpr("// bar"),
					},
					{
						Key:   ast.NewTextExpr("foo2"),
						Value: ast.NewTextExpr(`"bar2"`),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// foo"),
						},
						CommentExpr: ast.NewTextExpr("// bar"),
					},
					{
						Key: ast.NewTextExpr("foo3"),
						Value: ast.NewTextExpr(`"foo
			bar"`),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("/**foo*/"),
						},
						CommentExpr: ast.NewTextExpr("/**bar*/"),
					},
				},
			},
			ServiceApi: &ast.ServiceApi{
				ServiceToken: ast.NewTextExpr("service"),
				Name:         ast.NewTextExpr("foo-api"),
				Lbrace:       ast.NewTextExpr("{"),
				Rbrace:       ast.NewTextExpr("}"),
				ServiceRoute: []*ast.ServiceRoute{
					{
						AtDoc: &ast.AtDoc{
							AtDocToken: ast.NewTextExpr("@doc"),
							Lp:         ast.NewTextExpr("("),
							Rp:         ast.NewTextExpr(")"),
							LineDoc:    ast.NewTextExpr(`"foo"`),
						},
						AtHandler: &ast.AtHandler{
							AtHandlerToken: ast.NewTextExpr("@handler"),
							Name:           ast.NewTextExpr("foo"),
							DocExpr: []ast.Expr{
								ast.NewTextExpr("// foo"),
							},
							CommentExpr: ast.NewTextExpr("// bar"),
						},
						Route: &ast.Route{
							Method: ast.NewTextExpr("post"),
							Path:   ast.NewTextExpr("/foo"),
							Req: &ast.Body{
								Lp:   ast.NewTextExpr("("),
								Rp:   ast.NewTextExpr(")"),
								Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
							},
							ReturnToken: ast.NewTextExpr("returns"),
							Reply: &ast.Body{
								Lp:   ast.NewTextExpr("("),
								Rp:   ast.NewTextExpr(")"),
								Name: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
							},
							DocExpr: []ast.Expr{
								ast.NewTextExpr("// foo"),
							},
							CommentExpr: ast.NewTextExpr("// bar"),
						},
					},
				},
			},
		}))

		v, err = parser.Accept(fn, `
		service foo-api{
			@doc("foo")
			// foo/bar
			// foo
			@handler foo // bar
			// foo/bar
			// foo
			post /foo (Foo) returns (Bar) // bar
		}
		`)
		assert.Nil(t, err)
		service = v.(*ast.Service)
		assert.True(t, service.Equal(&ast.Service{
			ServiceApi: &ast.ServiceApi{
				ServiceToken: ast.NewTextExpr("service"),
				Name:         ast.NewTextExpr("foo-api"),
				Lbrace:       ast.NewTextExpr("{"),
				Rbrace:       ast.NewTextExpr("}"),
				ServiceRoute: []*ast.ServiceRoute{
					{
						AtDoc: &ast.AtDoc{
							AtDocToken: ast.NewTextExpr("@doc"),
							Lp:         ast.NewTextExpr("("),
							Rp:         ast.NewTextExpr(")"),
							LineDoc:    ast.NewTextExpr(`"foo"`),
						},
						AtHandler: &ast.AtHandler{
							AtHandlerToken: ast.NewTextExpr("@handler"),
							Name:           ast.NewTextExpr("foo"),
							DocExpr: []ast.Expr{
								ast.NewTextExpr("// foo"),
							},
							CommentExpr: ast.NewTextExpr("// bar"),
						},
						Route: &ast.Route{
							Method: ast.NewTextExpr("post"),
							Path:   ast.NewTextExpr("/foo"),
							Req: &ast.Body{
								Lp:   ast.NewTextExpr("("),
								Rp:   ast.NewTextExpr(")"),
								Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
							},
							ReturnToken: ast.NewTextExpr("returns"),
							Reply: &ast.Body{
								Lp:   ast.NewTextExpr("("),
								Rp:   ast.NewTextExpr(")"),
								Name: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
							},
							DocExpr: []ast.Expr{
								ast.NewTextExpr("// foo"),
							},
							CommentExpr: ast.NewTextExpr("// bar"),
						},
					},
				},
			},
		}))
	})
}
