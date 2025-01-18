package parser

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/assertx"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/ast"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/token"
)

//go:embed testdata/comment_test.api
var testCommentInput string

func TestParser_init(t *testing.T) {
	var testData = []string{
		"`",
		"@`",
		"syntax/**/`",
	}
	for _, val := range testData {
		p := New("test.api", val)
		val := p.init()
		assert.False(t, val)
	}
}

//go:embed testdata/test.api
var testInput string

func TestParser_Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) { // EXPERIMENTAL: just for testing output formatter.
		p := New("test.api", testInput)
		result := p.Parse()
		assert.NotNil(t, result)
		assert.True(t, p.hasNoErrors())
	})
	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			"foo bar",
			"@",
		}
		for _, val := range testData {
			p := New("test.api", val)
			p.Parse()
			assertx.ErrorOrigin(t, val, p.errors...)
		}
	})
}

func TestParser_Parse_Mode(t *testing.T) {

	t.Run("All", func(t *testing.T) {
		var testData = []string{
			`// foo`,
			`// bar`,
			`/*foo*/`,
			`/*bar*/`,
			`//baz`,
		}
		p := New("foo.api", testCommentInput)
		result := p.Parse()
		for idx, v := range testData {
			stmt := result.Stmts[idx]
			c, ok := stmt.(*ast.CommentStmt)
			assert.True(t, ok)
			assert.True(t, p.hasNoErrors())
			assert.Equal(t, v, c.Format(""))
		}
	})
}

func TestParser_Parse_syntaxStmt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []struct {
			input    string
			expected string
		}{
			{
				input:    `syntax = "v1"`,
				expected: `syntax = "v1"`,
			},
			{
				input:    `syntax = "foo"`,
				expected: `syntax = "foo"`,
			},
			{
				input:    `syntax= "bar"`,
				expected: `syntax = "bar"`,
			},
			{
				input:    ` syntax = "" `,
				expected: `syntax = ""`,
			},
		}
		for _, v := range testData {
			p := New("foo.aoi", v.input)
			result := p.Parse()
			assert.True(t, p.hasNoErrors())
			assert.Equal(t, v.expected, result.Stmts[0].Format(""))
		}
	})
	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`syntax`,
			`syntax = `,
			`syntax = ''`,
			`syntax = @`,
			`syntax = "v1`,
			`syntax == "v"`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/info_test.api
var infoTestAPI string

func TestParser_Parse_infoStmt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := map[string]string{
			"title":   `"type title here"`,
			"desc":    `"type desc here"`,
			"author":  `"type author here"`,
			"email":   `"type email here"`,
			"version": `"type version here"`,
		}
		p := New("foo.api", infoTestAPI)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		infoStmt, ok := stmt.(*ast.InfoStmt)
		assert.True(t, ok)
		for _, stmt := range infoStmt.Values {
			actual := stmt.Value.Token.Text
			expectedValue := expected[stmt.Key.Token.Text]
			assert.Equal(t, expectedValue, actual)
		}

	})

	t.Run("empty", func(t *testing.T) {
		p := New("foo.api", "info ()")
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		infoStmt, ok := stmt.(*ast.InfoStmt)
		assert.True(t, ok)
		assert.Equal(t, "info", infoStmt.Info.Token.Text)
		assert.Equal(t, "(", infoStmt.LParen.Token.Text)
		assert.Equal(t, ")", infoStmt.RParen.Token.Text)
		assert.Equal(t, 0, len(infoStmt.Values))
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`info`,
			`info(`,
			`info{`,
			`info(}`,
			`info( foo`,
			`info( foo:`,
			`info( foo:""`,
			`info( foo:"" bar`,
			`info( foo:"" bar:`,
			`info( foo:"" bar:""`,
			`info( foo:"`,
			`info foo:""`,
			`info( foo,""`,
			`info( foo-bar:"")`,
			`info(123:"")`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/import_literal_test.api
var testImportLiteral string

func TestParser_Parse_importLiteral(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []string{
			`""`,
			`"foo"`,
			`"bar"`,
		}
		p := New("foo.api", testImportLiteral)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		for idx, v := range testData {
			stmt := result.Stmts[idx]
			importLit, ok := stmt.(*ast.ImportLiteralStmt)
			assert.True(t, ok)
			assert.Equal(t, v, importLit.Value.Token.Text)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`import`,
			`import "`,
			`import "foo`,
			`import foo`,
			`import @`,
			`import $`,
			`import 好`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/import_group_test.api
var testImportGroup string

func TestParser_Parse_importGroup(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []string{
			`""`,
			`"foo"`,
			`"bar"`,
		}
		p := New("foo.api", testImportGroup)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		importGroup, ok := stmt.(*ast.ImportGroupStmt)
		assert.Equal(t, token.IDENT, importGroup.Import.Token.Type)
		assert.Equal(t, token.LPAREN, importGroup.LParen.Token.Type)
		assert.Equal(t, token.RPAREN, importGroup.RParen.Token.Type)
		for idx, v := range testData {
			assert.True(t, ok)
			assert.Equal(t, v, importGroup.Values[idx].Token.Text)
		}
	})

	t.Run("empty", func(t *testing.T) {
		p := New("foo.api", "import ()")
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		importGroup, ok := stmt.(*ast.ImportGroupStmt)
		assert.True(t, ok)
		assert.Equal(t, token.IDENT, importGroup.Import.Token.Type)
		assert.Equal(t, token.LPAREN, importGroup.LParen.Token.Type)
		assert.Equal(t, token.RPAREN, importGroup.RParen.Token.Type)
		assert.Equal(t, 0, len(importGroup.Values))
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`import`,
			`import (`,
			`import {`,
			`import ( "`,
			`import (} "`,
			`import ( ")`,
			`import ( ""`,
			`import ( "" foo)`,
			`import ( "" 好)`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/atserver_test.api
var atServerTestAPI string

func TestParser_Parse_atServerStmt(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var expectedData = map[string]string{
			"foo":        `bar`,
			"bar":        `baz`,
			"baz":        `foo`,
			"qux":        `/v1`,
			"quux":       `/v1/v2`,
			"middleware": `M1,M2`,
			"timeout1":   "1h",
			"timeout2":   "10m",
			"timeout3":   "10s",
			"timeout4":   "10ms",
			"timeout5":   "10µs",
			"timeout6":   "10ns",
			"timeout7":   "1h10m10s10ms10µs10ns",
			"maxBytes":   `1024`,
			"prefix":     "/v1",
			"prefix1":    "/v1/v2_test/v2-beta",
			"prefix2":    "v1/v2_test/v2-beta",
			"prefix3":    "v1/v2_",
			"prefix4":    "a-b-c",
			"summary":    `"test"`,
			"key":        `"bar"`,
		}

		p := New("foo.api", atServerTestAPI)
		result := p.ParseForUintTest()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		atServerStmt, ok := stmt.(*ast.AtServerStmt)
		assert.True(t, ok)
		for _, v := range atServerStmt.Values {
			expectedValue := expectedData[v.Key.Token.Text]
			assert.Equal(t, expectedValue, v.Value.Token.Text)
		}
	})

	t.Run("empty", func(t *testing.T) {
		p := New("foo.api", `@server()`)
		result := p.ParseForUintTest()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		atServerStmt, ok := stmt.(*ast.AtServerStmt)
		assert.True(t, ok)
		assert.Equal(t, token.AT_SERVER, atServerStmt.AtServer.Token.Type)
		assert.Equal(t, token.LPAREN, atServerStmt.LParen.Token.Type)
		assert.Equal(t, token.RPAREN, atServerStmt.RParen.Token.Type)
		assert.Equal(t, 0, len(atServerStmt.Values))
	})

	t.Run("invalidInSkipCommentMode", func(t *testing.T) {
		var testData = []string{
			`@server`,
			`@server{`,
			`@server(`,
			`@server(}`,
			`@server( //foo`,
			`@server( foo`,
			`@server( foo:`,
			`@server( foo:bar bar`,
			`@server( foo:bar bar,`,
			`@server( foo:bar bar: 123`,
			`@server( foo:bar bar: ""`,
			`@server( foo:bar bar: @`,
			`@server("":foo)`,
			`@server(foo:bar,baz)`,
			`@server(foo:/`,
			`@server(foo:/v`,
			`@server(foo:/v1/`,
			`@server(foo:/v1/v`,
			`@server(foo:/v1/v2`,
			`@server(foo: m1,`,
			`@server(foo: m1,)`,
			`@server(foo: v1/v2-)`,
			`@server(foo:"test")`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})

	t.Run("invalidWithNoSkipCommentMode", func(t *testing.T) {
		var testData = []string{
			`@server`,
			`@server //foo`,
			`@server /*foo*/`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.Error(t, p.errors...)
		}
	})
}

//go:embed testdata/athandler_test.api
var atHandlerTestAPI string

func TestParser_Parse_atHandler(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []string{
			`@handler foo`,
			`@handler foo1`,
			`@handler _bar`,
		}

		p := New("foo.api", atHandlerTestAPI)
		result := p.ParseForUintTest()
		assert.True(t, p.hasNoErrors())
		for idx, v := range testData {
			stmt := result.Stmts[idx]
			atHandlerStmt, ok := stmt.(*ast.AtHandlerStmt)
			assert.True(t, ok)
			assert.Equal(t, v, atHandlerStmt.Format(""))
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`@handler`,
			`@handler 1`,
			`@handler ""`,
			`@handler @`,
			`@handler $`,
			`@handler ()`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.ParseForUintTest()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/atdoc_literal_test.api
var atDocLiteralTestAPI string

func TestParser_Parse_atDocLiteral(t *testing.T) {
	t.Run("validLiteral", func(t *testing.T) {
		var testData = []string{
			`""`,
			`"foo"`,
			`"bar"`,
		}

		p := New("foo.api", atDocLiteralTestAPI)
		result := p.ParseForUintTest()
		assert.True(t, p.hasNoErrors())
		for idx, v := range testData {
			stmt := result.Stmts[idx]
			atDocLitStmt, ok := stmt.(*ast.AtDocLiteralStmt)
			assert.True(t, ok)
			assert.Equal(t, v, atDocLitStmt.Value.Token.Text)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`@doc`,
			`@doc "`,
			`@doc $`,
			`@doc 好`,
			`@doc |`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.ParseForUintTest()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/atdoc_group_test.api
var atDocGroupTestAPI string

func TestParser_Parse_atDocGroup(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = `@doc (
	foo: "foo"
	bar: "bar"
	baz: ""
)`

		p := New("foo.api", atDocGroupTestAPI)
		result := p.ParseForUintTest()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		atDocGroupStmt, _ := stmt.(*ast.AtDocGroupStmt)
		assert.Equal(t, testData, atDocGroupStmt.Format(""))
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`@doc{`,
			`@doc(`,
			`@doc(}`,
			`@doc( foo`,
			`@doc( foo:`,
			`@doc( foo: 123`,
			`@doc( foo: )`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.ParseForUintTest()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/service_test.api
var serviceTestAPI string

func TestParser_Parse_service(t *testing.T) {
	assertEqual := func(t *testing.T, expected, actual *ast.ServiceStmt) {
		if expected.AtServerStmt == nil {
			assert.Nil(t, actual.AtServerStmt)
		}
		assert.Equal(t, expected.Service.Token.Type, actual.Service.Token.Type)
		assert.Equal(t, expected.Service.Token.Text, actual.Service.Token.Text)
		assert.Equal(t, expected.Name.Format(""), actual.Name.Format(""))
		assert.Equal(t, expected.LBrace.Token.Type, actual.LBrace.Token.Type)
		assert.Equal(t, expected.RBrace.Token.Text, actual.RBrace.Token.Text)
		assert.Equal(t, len(expected.Routes), len(actual.Routes))
		for idx, v := range expected.Routes {
			actualItem := actual.Routes[idx]
			if v.AtDoc == nil {
				assert.Nil(t, actualItem.AtDoc)
			} else {
				assert.Equal(t, v.AtDoc.Format(""), actualItem.AtDoc.Format(""))
			}
			assert.Equal(t, v.AtHandler.Format(""), actualItem.AtHandler.Format(""))
			assert.Equal(t, v.Route.Format(""), actualItem.Route.Format(""))
		}
	}

	t.Run("valid", func(t *testing.T) {
		var testData = []*ast.ServiceStmt{
			{
				Service: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "service"}),
				Name: &ast.ServiceNameExpr{
					Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
				},
				LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
				RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
				Routes: []*ast.ServiceItemStmt{
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "root"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{Value: ast.NewTokenNode(token.Token{
								Type: token.PATH,
								Text: "/",
							})},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "bar"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{Value: ast.NewTokenNode(token.Token{
								Type: token.PATH,
								Text: "/ping",
							})},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "bar"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{Value: ast.NewTokenNode(token.Token{
								Type: token.PATH,
								Text: "/ping",
							})},
						},
					},
				},
			},
			{
				Service: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "service"}),
				Name: &ast.ServiceNameExpr{
					Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "bar"}),
				},
				LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
				RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
				Routes: []*ast.ServiceItemStmt{
					{
						AtDoc: &ast.AtDocLiteralStmt{
							AtDoc: ast.NewTokenNode(token.Token{Type: token.AT_DOC, Text: "@doc"}),
							Value: ast.NewTokenNode(token.Token{Type: token.STRING, Text: `"bar"`}),
						},
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "root"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtDoc: &ast.AtDocLiteralStmt{
							AtDoc: ast.NewTokenNode(token.Token{Type: token.AT_DOC, Text: "@doc"}),
							Value: ast.NewTokenNode(token.Token{Type: token.STRING, Text: `"bar"`}),
						},
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "root2"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/",
								}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtDoc: &ast.AtDocLiteralStmt{
							AtDoc: ast.NewTokenNode(token.Token{Type: token.AT_DOC, Text: "@doc"}),
							Value: ast.NewTokenNode(token.Token{Type: token.STRING, Text: `"bar"`}),
						},
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "root3"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtDoc: &ast.AtDocLiteralStmt{
							AtDoc: ast.NewTokenNode(token.Token{Type: token.AT_DOC, Text: "@doc"}),
							Value: ast.NewTokenNode(token.Token{Type: token.STRING, Text: `"bar"`}),
						},
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo/:bar",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtDoc: &ast.AtDocLiteralStmt{
							AtDoc: ast.NewTokenNode(token.Token{Type: token.AT_DOC, Text: "@doc"}),
							Value: ast.NewTokenNode(token.Token{Type: token.STRING, Text: `"bar"`}),
						},
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo/:bar",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo/:bar",
								}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "get"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo/:bar",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
				},
			},
			{
				Service: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "service"}),
				Name: &ast.ServiceNameExpr{
					Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "baz-api"}),
				},
				LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
				RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
				Routes: []*ast.ServiceItemStmt{
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "post"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo/:bar/foo-bar-baz",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Star:  ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "foo"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "post"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo/:bar/foo-bar-baz",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									Star:  ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
									Value: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "bar"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "post"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									LBrack: ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
									RBrack: ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
									Value:  ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									LBrack: ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
									RBrack: ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
									Star:   ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
									Value:  ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: ast.NewTokenNode(token.Token{Type: token.AT_HANDLER, Text: "@handler"}),
							Name:      ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "bar"}),
						},
						Route: &ast.RouteStmt{
							Method: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "post"}),
							Path: &ast.PathExpr{
								Value: ast.NewTokenNode(token.Token{
									Type: token.PATH,
									Text: "/foo",
								}),
							},
							Request: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									LBrack: ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
									RBrack: ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
									Value:  ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
							Returns: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "returns"}),
							Response: &ast.BodyStmt{
								LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
								Body: &ast.BodyExpr{
									LBrack: ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
									RBrack: ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
									Star:   ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
									Value:  ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
								},
								RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
							},
						},
					},
				},
			},
		}

		p := New("foo.api", serviceTestAPI)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		for idx, v := range testData {
			stmt := result.Stmts[idx]
			serviceStmt, ok := stmt.(*ast.ServiceStmt)
			assert.True(t, ok)
			assertEqual(t, v, serviceStmt)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`service`,
			`service foo`,
			`service -`,
			`service foo-`,
			`service foo-api`,
			`service foo(`,
			`service foo$`,
			`service foo好`,
			`service foo{`,
			`service foo{ @doc`,
			`service foo{ @doc $`,
			`service foo{ @doc ""`,
			`service foo{ @handler`,
			`service foo{ @handler foo`,
			`service foo{ @handler foo bar`,
			`service foo{ @handler foo get`,
			`service foo{ @handler foo get /`,
			`service foo{ @handler foo get \`,
			`service foo{ @handler foo get /:`,
			`service foo{ @handler foo get /::`,
			`service foo{ @handler foo get /:foo-`,
			`service foo{ @handler foo get /:foo--`,
			`service foo{ @handler foo get /:foo-bar/-`,
			`service foo{ @handler foo get /:foo-bar/baz`,
			`service foo{ @handler foo get /:foo-bar/baz (`,
			`service foo{ @handler foo get /:foo-bar/baz foo`,
			`service foo{ @handler foo get /:foo-bar/baz (foo`,
			`service foo{ @handler foo get /:foo-bar/baz (foo)`,
			`service foo{ @handler foo get /:foo-bar/baz (foo) return`,
			`service foo{ @handler foo get /:foo-bar/baz (foo) returns `,
			`service foo{ @handler foo get /:foo-bar/baz (foo) returns (`,
			`service foo{ @handler foo get /:foo-bar/baz (foo) returns {`,
			`service foo{ @handler foo get /:foo-bar/baz (foo) returns (bar`,
			`service foo{ @handler foo get /:foo-bar/baz (foo) returns (bar}`,
			`service foo{ @handler foo get /:foo-bar/baz (foo) returns (bar)`,
			`service foo{ @handler foo get /:foo-bar/baz ([`,
			`service foo{ @handler foo get /:foo-bar/baz ([@`,
			`service foo{ @handler foo get /:foo-bar/baz ([]`,
			`service foo{ @handler foo get /:foo-bar/baz ([]*`,
			`service foo{ @handler foo get /:foo-bar/baz ([]*Foo`,
			`service foo{ @handler foo get /:foo-bar/baz (*`,
			`service foo{ @handler foo get /:foo-bar/baz (*Foo`,
			`service foo{ @handler foo get /:foo-bar/baz returns`,
			`service foo{ @handler foo get /:foo-bar/baz returns (`,
			`service foo{ @handler foo get /:foo-bar/baz returns ([`,
			`service foo{ @handler foo get /:foo-bar/baz returns ([]`,
			`service foo{ @handler foo get /:foo-bar/baz returns ([]*`,
			`service foo{ @handler foo get /:foo-bar/baz returns ([]*Foo`,
			`service foo{ @handler foo get /:foo-bar/baz returns (*`,
			`service foo{ @handler foo get /:foo-bar/baz returns (*Foo`,
			`service foo{ @handler foo get /ping (Foo) returns (Bar)]`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})

	t.Run("invalidBeginWithAtServer", func(t *testing.T) {
		var testData = []string{
			`@server(`,
			`@server() service`,
			`@server() foo`,
			`@server() service fo`,
		}
		for _, v := range testData {
			p := New("foo.api", v)
			p.init()
			_ = p.parseService()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

func TestParser_Parse_pathItem(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var testData = []struct {
			input    string
			expected string
		}{
			{input: "foo", expected: "foo"},
			{input: "foo2", expected: "foo2"},
			{input: "foo-bar", expected: "foo-bar"},
			{input: "foo-bar2", expected: "foo-bar2"},
			{input: "foo-bar-baz", expected: "foo-bar-baz"},
			{input: "_foo-bar-baz", expected: "_foo-bar-baz"},
			{input: "_foo_bar-baz", expected: "_foo_bar-baz"},
			{input: "_foo_bar_baz", expected: "_foo_bar_baz"},
			{input: "_foo_bar_baz", expected: "_foo_bar_baz"},
			{input: "foo/", expected: "foo"},
			{input: "foo(", expected: "foo"},
			{input: "foo returns", expected: "foo"},
			{input: "foo @doc", expected: "foo"},
			{input: "foo @handler", expected: "foo"},
			{input: "foo }", expected: "foo"},
			{input: "1", expected: "1"},
			{input: "11", expected: "11"},
		}
		for _, v := range testData {
			p := New("foo.api", v.input)
			ok := p.nextToken()
			assert.True(t, ok)
			tokens := p.parsePathItem()
			var expected []string
			for _, tok := range tokens {
				expected = append(expected, tok.Text)
			}
			assert.Equal(t, strings.Join(expected, ""), v.expected)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			"-foo",
			"foo-",
			"foo-2",
			"2-2",
			"foo-bar-123",
			"foo-bar-$",
			"foo-bar-好",
			"foo-bar@",
			"foo-barの",
		}
		for _, v := range testData {
			p := New("foo.api", v)
			ok := p.nextToken()
			assert.True(t, ok)
			p.parsePathItem()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

func TestParser_Parse_parseTypeStmt(t *testing.T) {
	assertEqual := func(t *testing.T, expected, actual ast.Stmt) {
		if expected == nil {
			assert.Nil(t, actual)
			return
		}
		assert.Equal(t, expected.Format(""), actual.Format(""))
	}
	t.Run("parseTypeLiteralStmt", func(t *testing.T) {
		var testData = []struct {
			input    string
			expected ast.TypeStmt
		}{
			{
				input: "type Int int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Int"}),
						DataType: &ast.BaseDataType{
							Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"}),
						},
					},
				},
			},
			{
				input: "type Int interface{}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Int"}),
						DataType: &ast.InterfaceDataType{
							Interface: ast.NewTokenNode(token.Token{Type: token.ANY, Text: "interface{}"}),
						},
					},
				},
			},
			{
				input: "type Int any",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Int"}),
						DataType: &ast.AnyDataType{
							Any: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "any"}),
						},
					},
				},
			},
			{
				input: "type Int = int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name:   ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Int"}),
						Assign: ast.NewTokenNode(token.Token{Type: token.ASSIGN, Text: "="}),
						DataType: &ast.BaseDataType{
							Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"}),
						},
					},
				},
			},
			{
				input: "type Array [2]int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Array"}),
						DataType: &ast.ArrayDataType{
							LBrack:   ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
							Length:   ast.NewTokenNode(token.Token{Type: token.INT, Text: "2"}),
							RBrack:   ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
							DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
				},
			},
			{
				input: "type Array [...]int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Array"}),
						DataType: &ast.ArrayDataType{
							LBrack:   ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
							Length:   ast.NewTokenNode(token.Token{Type: token.ELLIPSIS, Text: "..."}),
							RBrack:   ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
							DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
				},
			},
			{
				input: "type Map map[string]int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Map"}),
						DataType: &ast.MapDataType{
							Map:    ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "map"}),
							LBrack: ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
							Key:    &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "string"})},
							RBrack: ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
							Value:  &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
				},
			},
			{
				input: "type Pointer *int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Pointer"}),
						DataType: &ast.PointerDataType{
							Star:     ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
							DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
				},
			},
			{
				input: "type Slice []int",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Slice"}),
						DataType: &ast.SliceDataType{
							LBrack:   ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
							RBrack:   ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
							DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
				},
			},
			{
				input: "type Foo {}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
			{
				input: "type Foo {Bar\n*Baz}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							Elements: ast.ElemExprList{
								{
									DataType: &ast.BaseDataType{
										Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
									},
								},
								{
									DataType: &ast.PointerDataType{
										Star: ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
										DataType: &ast.BaseDataType{
											Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Baz"}),
										},
									},
								},
							},
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
			{
				input: "type Foo {Bar `json:\"bar\"`\n*Baz `json:\"baz\"`}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							Elements: ast.ElemExprList{
								{
									DataType: &ast.BaseDataType{
										Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Bar"}),
									},
									Tag: ast.NewTokenNode(token.Token{
										Type: token.RAW_STRING,
										Text: "`json:\"bar\"`",
									}),
								},
								{
									DataType: &ast.PointerDataType{
										Star: ast.NewTokenNode(token.Token{Type: token.MUL, Text: "*"}),
										DataType: &ast.BaseDataType{
											Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Baz"}),
										},
									},
									Tag: ast.NewTokenNode(token.Token{
										Type: token.RAW_STRING,
										Text: "`json:\"baz\"`",
									}),
								},
							},
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
			{
				input: "type Foo {Name string}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							Elements: ast.ElemExprList{
								{
									Name: []*ast.TokenNode{ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Name"})},
									DataType: &ast.BaseDataType{
										Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "string"}),
									},
								},
							},
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
			{
				input: "type Foo {Name,Desc string}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							Elements: ast.ElemExprList{
								{
									Name: []*ast.TokenNode{
										{Token: token.Token{Type: token.IDENT, Text: "Name"}},
										{Token: token.Token{Type: token.IDENT, Text: "Desc"}},
									},
									DataType: &ast.BaseDataType{
										Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "string"}),
									},
								},
							},
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
			{
				input: "type Foo {Name string\n Age int `json:\"age\"`}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							Elements: ast.ElemExprList{
								{
									Name: []*ast.TokenNode{
										{Token: token.Token{Type: token.IDENT, Text: "Name"}},
									},
									DataType: &ast.BaseDataType{
										Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "string"}),
									},
								},
								{
									Name: []*ast.TokenNode{
										{Token: token.Token{Type: token.IDENT, Text: "Age"}},
									},
									DataType: &ast.BaseDataType{
										Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"}),
									},
									Tag: ast.NewTokenNode(token.Token{Type: token.RAW_STRING, Text: "`json:\"age\"`"}),
								},
							},
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
			{
				input: "type Foo {Bar {Name string}}",
				expected: &ast.TypeLiteralStmt{
					Type: ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					Expr: &ast.TypeExpr{
						Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Foo"}),
						DataType: &ast.StructDataType{
							LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
							Elements: ast.ElemExprList{
								{
									Name: []*ast.TokenNode{
										{Token: token.Token{Type: token.IDENT, Text: "Bar"}},
									},
									DataType: &ast.StructDataType{
										LBrace: ast.NewTokenNode(token.Token{Type: token.LBRACE, Text: "{"}),
										Elements: ast.ElemExprList{
											{
												Name: []*ast.TokenNode{
													{Token: token.Token{Type: token.IDENT, Text: "Name"}},
												},
												DataType: &ast.BaseDataType{
													Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "string"}),
												},
											},
										},
										RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
									},
								},
							},
							RBrace: ast.NewTokenNode(token.Token{Type: token.RBRACE, Text: "}"}),
						},
					},
				},
			},
		}
		for _, val := range testData {
			p := New("test.api", val.input)
			result := p.Parse()
			assert.True(t, p.hasNoErrors())
			assert.Equal(t, 1, len(result.Stmts))
			one := result.Stmts[0]
			assertEqual(t, val.expected, one)
		}
	})
	t.Run("parseTypeGroupStmt", func(t *testing.T) {
		var testData = []struct {
			input    string
			expected ast.TypeStmt
		}{
			{
				input: "type (Int int)",
				expected: &ast.TypeGroupStmt{
					Type:   ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
					ExprList: []*ast.TypeExpr{
						{
							Name:     ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Int"}),
							DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
					RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
				},
			},
			{
				input: "type (Int = int)",
				expected: &ast.TypeGroupStmt{
					Type:   ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
					ExprList: []*ast.TypeExpr{
						{
							Name:     ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Int"}),
							Assign:   ast.NewTokenNode(token.Token{Type: token.ASSIGN, Text: "="}),
							DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
						},
					},
					RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
				},
			},
			{
				input: "type (Array [2]int)",
				expected: &ast.TypeGroupStmt{
					Type:   ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
					ExprList: []*ast.TypeExpr{
						{
							Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Array"}),
							DataType: &ast.ArrayDataType{
								LBrack:   ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
								Length:   ast.NewTokenNode(token.Token{Type: token.INT, Text: "2"}),
								RBrack:   ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
								DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
							},
						},
					},
					RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
				},
			},
			{
				input: "type (Array [...]int)",
				expected: &ast.TypeGroupStmt{
					Type:   ast.NewTokenNode(token.Token{Type: token.TYPE, Text: "type"}),
					LParen: ast.NewTokenNode(token.Token{Type: token.LPAREN, Text: "("}),
					ExprList: []*ast.TypeExpr{
						{
							Name: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "Array"}),
							DataType: &ast.ArrayDataType{
								LBrack:   ast.NewTokenNode(token.Token{Type: token.LBRACK, Text: "["}),
								Length:   ast.NewTokenNode(token.Token{Type: token.ELLIPSIS, Text: "..."}),
								RBrack:   ast.NewTokenNode(token.Token{Type: token.RBRACK, Text: "]"}),
								DataType: &ast.BaseDataType{Base: ast.NewTokenNode(token.Token{Type: token.IDENT, Text: "int"})},
							},
						},
					},
					RParen: ast.NewTokenNode(token.Token{Type: token.RPAREN, Text: ")"}),
				},
			},
		}
		for _, val := range testData {
			p := New("test.api", val.input)
			result := p.Parse()
			assert.True(t, p.hasNoErrors())
			assert.Equal(t, 1, len(result.Stmts))
			one := result.Stmts[0]
			assertEqual(t, val.expected, one)
		}
	})
}

func TestParser_Parse_parseTypeStmt_invalid(t *testing.T) {
	var testData = []string{
		/**************** type literal stmt ****************/
		"type",
		"type @",
		"type Foo",
		"type Foo = ",
		"type Foo = [",
		"type Foo = []",
		"type Foo = [2",
		"type Foo = [2]",
		"type Foo = [...",
		"type Foo = [...]",
		"type Foo map",
		"type Foo map[",
		"type Foo map[]",
		"type Foo map[string",
		"type Foo map[123",
		"type Foo map[string]",
		"type Foo map[string]@",
		"type Foo interface",
		"type Foo interface{",
		"type Foo *",
		"type Foo *123",
		"type Foo *go",
		"type Foo {",
		"type Foo { Foo ",
		"type Foo { Foo int",
		"type Foo { Foo int `",
		"type Foo { Foo int ``",
		"type Foo { Foo,",
		"type Foo { Foo@",
		"type Foo { Foo,Bar",
		"type Foo { Foo,Bar int",
		"type Foo { Foo,Bar int `baz`",
		"type Foo { Foo,Bar int `baz`)",
		"type Foo { Foo,Bar int `baz`@",
		"type Foo *",
		"type Foo *{",
		"type Foo *[",
		"type Foo *[]",
		"type Foo *map",
		"type Foo *map[",
		"type Foo *map[int",
		"type Foo *map[int]123",
		"type Foo *map[int]@",
		"type Foo *好",

		/**************** type group stmt ****************/
		"type (@",
		"type (Foo",
		"type (Foo = ",
		"type (Foo = [",
		"type (Foo = []",
		"type (Foo = [2",
		"type (Foo = [2]",
		"type (Foo = [...",
		"type (Foo = [...]",
		"type (Foo map",
		"type (Foo map[",
		"type (Foo map[]",
		"type (Foo map[string",
		"type (Foo map[123",
		"type (Foo map[string]",
		"type (Foo map[string]@",
		"type (Foo interface",
		"type (Foo interface{",
		"type (Foo *",
		"type (Foo *123",
		"type (Foo *go",
		"type (Foo {",
		"type (Foo { Foo ",
		"type (Foo { Foo int",
		"type (Foo { Foo int `",
		"type (Foo { Foo int ``",
		"type (Foo { Foo,",
		"type (Foo { Foo@",
		"type (Foo { Foo,Bar",
		"type (Foo { Foo,Bar int",
		"type (Foo { Foo,Bar int `baz`",
		"type (Foo { Foo,Bar int `baz`)",
		"type (Foo { Foo,Bar int `baz`@",
		"type (Foo *",
		"type (Foo *{",
		"type (Foo *[",
		"type (Foo *[]",
		"type (Foo *map",
		"type (Foo *map[",
		"type (Foo *map[int",
		"type (Foo *map[int]123",
		"type (Foo *map[int]@",
		"type (Foo *好",
		"type (Foo)",
		"type (Foo int\nBar)",
		"type (Foo int\nBar string `)",
		"type (())",
		"type go int",
		"type  A func",
		"type  A {map[string]int}",
		"type  A {Name \n string}",
	}

	for _, v := range testData {
		p := New("test.api", v)
		p.Parse()
		assertx.ErrorOrigin(t, v, p.errors...)
	}
}
