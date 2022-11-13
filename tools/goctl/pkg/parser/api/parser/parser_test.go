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
		p := New("test.api", val, SkipComment)
		val := p.init()
		assert.False(t, val)
	}
}

//go:embed testdata/test.api
var testInput string

func TestParser_Parse(t *testing.T) {
	t.Run("valid", func(t *testing.T) { // EXPERIMENTAL: just for testing output formatter.
		p := New("test.api", testInput, SkipComment)
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
			p := New("test.api", val, SkipComment)
			p.Parse()
			assertx.ErrorOrigin(t, val, p.errors...)
		}
	})
}

func TestParser_Parse_Mode(t *testing.T) {
	t.Run("SkipComment", func(t *testing.T) {
		p := New("foo.api", testCommentInput, SkipComment)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		assert.Equal(t, 0, len(result.Stmts))
	})

	//t.Run("All", func(t *testing.T) {
	//	var testData = []string{
	//		`// foo`,
	//		`// bar`,
	//		`/*foo*/`,
	//		`/*bar*/`,
	//		`//baz`,
	//	}
	//	p := New("foo.api", testCommentInput, All)
	//	result := p.Parse()
	//	for idx, v := range testData {
	//		stmt := result.Stmts[idx]
	//		c, ok := stmt.(*ast.CommentStmt)
	//		assert.True(t, ok)
	//		assert.True(t, p.hasNoErrors())
	//		assert.Equal(t, v, c.Format(""))
	//	}
	//})
}

func TestParser_Parse_syntaxStmt(t *testing.T) {
	//t.Run("valid", func(t *testing.T) {
	//	var testData = []struct {
	//		input    string
	//		expected string
	//	}{
	//		{
	//			input:    `syntax = "v1"`,
	//			expected: `syntax = "v1"`,
	//		},
	//		{
	//			input:    `syntax = "foo"`,
	//			expected: `syntax = "foo"`,
	//		},
	//		{
	//			input:    `syntax= "bar"`,
	//			expected: `syntax = "bar"`,
	//		},
	//		{
	//			input:    ` syntax = "" `,
	//			expected: `syntax = ""`,
	//		},
	//	}
	//	for _, v := range testData {
	//		p := New("foo.aoi", v.input, SkipComment)
	//		result := p.Parse()
	//		assert.True(t, p.hasNoErrors())
	//		assert.Equal(t, v.expected, result.Stmts[0].Format(""))
	//	}
	//})
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
			p := New("foo.api", v, SkipComment)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/info_test.api
var infoTestAPI string

func TestParser_Parse_infoStmt(t *testing.T) {
	//t.Run("valid", func(t *testing.T) {
	//	var testData = []string{
	//		`title: "type title here"`,
	//		`desc: "type desc here"`,
	//		`author: "type author here"`,
	//		`email: "type email here"`,
	//		`version: "type version here"`,
	//	}
	//	p := New("foo.api", infoTestAPI, SkipComment)
	//	result := p.Parse()
	//	assert.True(t, p.hasNoErrors())
	//	stmt := result.Stmts[0]
	//	infoStmt, ok := stmt.(*ast.InfoStmt)
	//	assert.True(t, ok)
	//	for idx, v := range testData {
	//		assert.Equal(t, v, infoStmt.Values[idx].Format(""))
	//	}
	//})

	//t.Run("empty", func(t *testing.T) {
	//	p := New("foo.api", "info ()", SkipComment)
	//	result := p.Parse()
	//	assert.True(t, p.hasNoErrors())
	//	stmt := result.Stmts[0]
	//	infoStmt, ok := stmt.(*ast.InfoStmt)
	//	assert.True(t, ok)
	//	assert.Equal(t, "info ()", infoStmt.Format(""))
	//})

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
			p := New("foo.api", v, SkipComment)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/import_literal_test.api
var testImportLiteral string

func TestParser_Parse_importLiteral(t *testing.T) {
	//t.Run("valid", func(t *testing.T) {
	//	var testData = []string{
	//		`import ""`,
	//		`import "foo"`,
	//		`import "bar"`,
	//	}
	//	p := New("foo.api", testImportLiteral, SkipComment)
	//	result := p.Parse()
	//	assert.True(t, p.hasNoErrors())
	//	for idx, v := range testData {
	//		stmt := result.Stmts[idx]
	//		importLit, ok := stmt.(*ast.ImportLiteralStmt)
	//		assert.True(t, ok)
	//		assert.Equal(t, v, importLit.Format(""))
	//	}
	//})
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
			p := New("foo.api", v, SkipComment)
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
		p := New("foo.api", testImportGroup, SkipComment)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		importGroup, ok := stmt.(*ast.ImportGroupStmt)
		assert.Equal(t, token.IMPORT, importGroup.Import.Type)
		assert.Equal(t, token.LPAREN, importGroup.LParen.Type)
		assert.Equal(t, token.RPAREN, importGroup.RParen.Type)
		for idx, v := range testData {
			assert.True(t, ok)
			assert.Equal(t, v, importGroup.Values[idx].Text)
		}
	})

	t.Run("empty", func(t *testing.T) {
		p := New("foo.api", "import ()", SkipComment)
		result := p.Parse()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		importGroup, ok := stmt.(*ast.ImportGroupStmt)
		assert.True(t, ok)
		assert.Equal(t, token.IMPORT, importGroup.Import.Type)
		assert.Equal(t, token.LPAREN, importGroup.LParen.Type)
		assert.Equal(t, token.RPAREN, importGroup.RParen.Type)
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
			p := New("foo.api", v, SkipComment)
			_ = p.Parse()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/atserver_test.api
var atServerTestAPI string

func TestParser_Parse_atServerStmt(t *testing.T) {
	//t.Run("valid", func(t *testing.T) {
	//	var testData = []string{
	//		`foo: bar`,
	//		`bar: baz`,
	//		`baz: foo`,
	//		`qux: /v1`,
	//		`quux: /v1/v2`,
	//	}
	//	p := New("foo.api", atServerTestAPI, SkipComment)
	//	result := p.parseForUintTest()
	//	assert.True(t, p.hasNoErrors())
	//	stmt := result.Stmts[0]
	//	atServerStmt, ok := stmt.(*ast.AtServerStmt)
	//	assert.True(t, ok)
	//	for idx, v := range testData {
	//		assert.Equal(t, v, atServerStmt.Values[idx].Format(""))
	//	}
	//})

	t.Run("empty", func(t *testing.T) {
		p := New("foo.api", `@server()`, SkipComment)
		result := p.parseForUintTest()
		assert.True(t, p.hasNoErrors())
		stmt := result.Stmts[0]
		atServerStmt, ok := stmt.(*ast.AtServerStmt)
		assert.True(t, ok)
		assert.Equal(t, token.AT_SERVER, atServerStmt.AtServer.Type)
		assert.Equal(t, token.LPAREN, atServerStmt.LParen.Type)
		assert.Equal(t, token.RPAREN, atServerStmt.RParen.Type)
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
		}
		for _, v := range testData {
			p := New("foo.api", v, SkipComment)
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
			p := New("foo.api", v, All)
			_ = p.Parse()
			assertx.Error(t, p.errors...)
		}
	})
}

//go:embed testdata/athandler_test.api
var atHandlerTestAPI string

func TestParser_Parse_atHandler(t *testing.T) {
	//t.Run("valid", func(t *testing.T) {
	//	var testData = []string{
	//		`@handler foo`,
	//		`@handler foo1`,
	//		`@handler _bar`,
	//	}
	//
	//	p := New("foo.api", atHandlerTestAPI, SkipComment)
	//	result := p.parseForUintTest()
	//	assert.True(t, p.hasNoErrors())
	//	for idx, v := range testData {
	//		stmt := result.Stmts[idx]
	//		atHandlerStmt, ok := stmt.(*ast.AtHandlerStmt)
	//		assert.True(t, ok)
	//		assert.Equal(t, v, atHandlerStmt.Format(""))
	//	}
	//})

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
			p := New("foo.api", v, SkipComment)
			_ = p.parseForUintTest()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/atdoc_literal_test.api
var atDocLiteralTestAPI string

func TestParser_Parse_atDocLiteral(t *testing.T) {
	//t.Run("validLiteral", func(t *testing.T) {
	//	var testData = []string{
	//		`@doc ""`,
	//		`@doc "foo"`,
	//		`@doc "bar"`,
	//	}
	//
	//	p := New("foo.api", atDocLiteralTestAPI, SkipComment)
	//	result := p.parseForUintTest()
	//	assert.True(t, p.hasNoErrors())
	//	for idx, v := range testData {
	//		stmt := result.Stmts[idx]
	//		atDocLitStmt, ok := stmt.(*ast.AtDocLiteralStmt)
	//		assert.True(t, ok)
	//		assert.Equal(t, v, atDocLitStmt.Format(""))
	//	}
	//})

	t.Run("invalid", func(t *testing.T) {
		var testData = []string{
			`@doc`,
			`@doc "`,
			`@doc $`,
			`@doc 好`,
			`@doc |`,
		}
		for _, v := range testData {
			p := New("foo.api", v, SkipComment)
			_ = p.parseForUintTest()
			assertx.ErrorOrigin(t, v, p.errors...)
		}
	})
}

//go:embed testdata/atdoc_group_test.api
var atDocGroupTestAPI string

func TestParser_Parse_atDocGroup(t *testing.T) {
	//t.Run("valid", func(t *testing.T) {
	//	var testData = []string{
	//		`foo: "foo"`,
	//		`bar: "bar"`,
	//		`baz: ""`,
	//	}
	//
	//	p := New("foo.api", atDocGroupTestAPI, SkipComment)
	//	result := p.parseForUintTest()
	//	assert.True(t, p.hasNoErrors())
	//	stmt := result.Stmts[0]
	//	atDocLitStmt, ok := stmt.(*ast.AtDocGroupStmt)
	//	for idx, v := range testData {
	//		assert.True(t, ok)
	//		assert.Equal(t, v, atDocLitStmt.Values[idx].Format(""))
	//	}
	//})

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
			p := New("foo.api", v, SkipComment)
			_ = p.parseForUintTest()
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
		assert.Equal(t, expected.Service.Type, actual.Service.Type)
		assert.Equal(t, expected.Service.Text, actual.Service.Text)
		//assert.Equal(t, expected.Name.Format(""), actual.Name.Format(""))
		assert.Equal(t, expected.LBrace.Type, actual.LBrace.Type)
		assert.Equal(t, expected.RBrace.Text, actual.RBrace.Text)
		assert.Equal(t, len(expected.Routes), len(actual.Routes))
		for idx, v := range expected.Routes {
			actualItem := actual.Routes[idx]
			if v.AtDoc == nil {
				assert.Nil(t, actualItem.AtDoc)
			} else {
				//assert.Equal(t, v.AtDoc.Format(""), actualItem.AtDoc.Format(""))
			}
			//assert.Equal(t, v.AtHandler.Format(""), actualItem.AtHandler.Format(""))
			//assert.Equal(t, v.Route.Format(""), actualItem.Route.Format(""))
		}
	}

	t.Run("valid", func(t *testing.T) {
		var testData = []*ast.ServiceStmt{
			{
				Service: token.Token{Type: token.IDENT, Text: "service"},
				Name: &ast.ServiceNameExpr{
					ID: token.Token{Type: token.IDENT, Text: "foo"},
				},
				LBrace: token.Token{Type: token.LBRACE, Text: "{"},
				RBrace: token.Token{Type: token.RBRACE, Text: "}"},
				Routes: []*ast.ServiceItemStmt{
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: token.Token{Type: token.AT_HANDLER, Text: "@handler"},
							Name:      token.Token{Type: token.IDENT, Text: "bar"},
						},
						Route: &ast.RouteStmt{
							Method: token.Token{Type: token.IDENT, Text: "get"},
							Path: &ast.PathExpr{Values: []token.Token{
								{Type: token.QUO, Text: "/"},
								{Type: token.IDENT, Text: "ping"},
							}},
						},
					},
				},
			},
			{
				Service: token.Token{Type: token.IDENT, Text: "service"},
				Name: &ast.ServiceNameExpr{
					ID: token.Token{Type: token.IDENT, Text: "bar"},
				},
				LBrace: token.Token{Type: token.LBRACE, Text: "{"},
				RBrace: token.Token{Type: token.RBRACE, Text: "}"},
				Routes: []*ast.ServiceItemStmt{
					{
						AtDoc: &ast.AtDocLiteralStmt{
							AtDoc: token.Token{Type: token.AT_DOC, Text: "@doc"},
							Value: token.Token{Type: token.STRING, Text: `"bar"`},
						},
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: token.Token{Type: token.AT_HANDLER, Text: "@handler"},
							Name:      token.Token{Type: token.IDENT, Text: "foo"},
						},
						Route: &ast.RouteStmt{
							Method: token.Token{Type: token.IDENT, Text: "get"},
							Path: &ast.PathExpr{Values: []token.Token{
								{Type: token.QUO, Text: "/"},
								{Type: token.IDENT, Text: "foo"},
								{Type: token.QUO, Text: "/"},
								{Type: token.COLON, Text: ":"},
								{Type: token.IDENT, Text: "bar"},
							}},
							Request: &ast.BodyStmt{
								LParen: token.Token{Type: token.LPAREN, Text: "("},
								Body: &ast.BodyExpr{
									Value: token.Token{Type: token.IDENT, Text: "Foo"},
								},
								RParen: token.Token{Type: token.RPAREN, Text: ")"},
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: token.Token{Type: token.AT_HANDLER, Text: "@handler"},
							Name:      token.Token{Type: token.IDENT, Text: "foo"},
						},
						Route: &ast.RouteStmt{
							Method: token.Token{Type: token.IDENT, Text: "get"},
							Path: &ast.PathExpr{Values: []token.Token{
								{Type: token.QUO, Text: "/"},
								{Type: token.IDENT, Text: "foo"},
								{Type: token.QUO, Text: "/"},
								{Type: token.COLON, Text: ":"},
								{Type: token.IDENT, Text: "bar"},
							}},
							Returns: token.Token{Type: token.IDENT, Text: "returns"},
							Response: &ast.BodyStmt{
								LParen: token.Token{Type: token.LPAREN, Text: "("},
								Body: &ast.BodyExpr{
									Value: token.Token{Type: token.IDENT, Text: "Foo"},
								},
								RParen: token.Token{Type: token.RPAREN, Text: ")"},
							},
						},
					},
				},
			},
			{
				Service: token.Token{Type: token.IDENT, Text: "service"},
				Name: &ast.ServiceNameExpr{
					ID:     token.Token{Type: token.IDENT, Text: "baz"},
					Joiner: token.Token{Type: token.SUB, Text: "-"},
					API:    token.Token{Type: token.IDENT, Text: "api"},
				},
				LBrace: token.Token{Type: token.LBRACE, Text: "{"},
				RBrace: token.Token{Type: token.RBRACE, Text: "}"},
				Routes: []*ast.ServiceItemStmt{
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: token.Token{Type: token.AT_HANDLER, Text: "@handler"},
							Name:      token.Token{Type: token.IDENT, Text: "foo"},
						},
						Route: &ast.RouteStmt{
							Method: token.Token{Type: token.IDENT, Text: "post"},
							Path: &ast.PathExpr{Values: []token.Token{
								{Type: token.QUO, Text: "/"},
								{Type: token.IDENT, Text: "foo"},
								{Type: token.QUO, Text: "/"},
								{Type: token.COLON, Text: ":"},
								{Type: token.IDENT, Text: "bar"},
								{Type: token.QUO, Text: "/"},
								{Type: token.IDENT, Text: "foo-bar-baz"},
							}},
							Request: &ast.BodyStmt{
								LParen: token.Token{Type: token.LPAREN, Text: "("},
								Body: &ast.BodyExpr{
									Value: token.Token{Type: token.IDENT, Text: "Foo"},
								},
								RParen: token.Token{Type: token.RPAREN, Text: ")"},
							},
							Returns: token.Token{Type: token.IDENT, Text: "returns"},
							Response: &ast.BodyStmt{
								LParen: token.Token{Type: token.LPAREN, Text: "("},
								Body: &ast.BodyExpr{
									Star:  token.Token{Type: token.MUL, Text: "*"},
									Value: token.Token{Type: token.IDENT, Text: "Bar"},
								},
								RParen: token.Token{Type: token.RPAREN, Text: ")"},
							},
						},
					},
					{
						AtHandler: &ast.AtHandlerStmt{
							AtHandler: token.Token{Type: token.AT_HANDLER, Text: "@handler"},
							Name:      token.Token{Type: token.IDENT, Text: "bar"},
						},
						Route: &ast.RouteStmt{
							Method: token.Token{Type: token.IDENT, Text: "post"},
							Path: &ast.PathExpr{Values: []token.Token{
								{Type: token.QUO, Text: "/"},
								{Type: token.IDENT, Text: "foo"},
							}},
							Request: &ast.BodyStmt{
								LParen: token.Token{Type: token.LPAREN, Text: "("},
								Body: &ast.BodyExpr{
									LBrack: token.Token{Type: token.LBRACK, Text: "["},
									RBrack: token.Token{Type: token.RBRACK, Text: "]"},
									Value:  token.Token{Type: token.IDENT, Text: "Foo"},
								},
								RParen: token.Token{Type: token.RPAREN, Text: ")"},
							},
							Returns: token.Token{Type: token.IDENT, Text: "returns"},
							Response: &ast.BodyStmt{
								LParen: token.Token{Type: token.LPAREN, Text: "("},
								Body: &ast.BodyExpr{
									LBrack: token.Token{Type: token.LBRACK, Text: "["},
									RBrack: token.Token{Type: token.RBRACK, Text: "]"},
									Star:   token.Token{Type: token.MUL, Text: "*"},
									Value:  token.Token{Type: token.IDENT, Text: "Bar"},
								},
								RParen: token.Token{Type: token.RPAREN, Text: ")"},
							},
						},
					},
				},
			},
		}

		p := New("foo.api", serviceTestAPI, SkipComment)
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
			p := New("foo.api", v, SkipComment)
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
			p := New("foo.api", v, SkipComment)
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
		}
		for _, v := range testData {
			p := New("foo.api", v.input, SkipComment)
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
			"foo-bar-123",
			"foo-bar-$",
			"foo-bar-好",
			"foo-bar@",
			"foo-barの",
		}
		for _, v := range testData {
			p := New("foo.api", v, SkipComment)
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
		//assert.Equal(t, expected.Format(""), actual.Format(""))
	}
	t.Run("parseTypeLiteralStmt", func(t *testing.T) {
		var testData = []struct {
			input    string
			expected ast.TypeStmt
		}{
			{
				input: "type Int int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Int"},
						DataType: &ast.BaseDataType{
							Base: token.Token{Type: token.IDENT, Text: "int"},
						},
					},
				},
			},
			{
				input: "type Int interface{}",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Int"},
						DataType: &ast.InterfaceDataType{
							Interface: token.Token{Type: token.ANY, Text: "interface{}"},
						},
					},
				},
			},
			{
				input: "type Int any",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Int"},
						DataType: &ast.AnyDataType{
							Any: token.Token{Type: token.IDENT, Text: "any"},
						},
					},
				},
			},
			{
				input: "type Int = int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name:   token.Token{Type: token.IDENT, Text: "Int"},
						Assign: token.Token{Type: token.ASSIGN, Text: "="},
						DataType: &ast.BaseDataType{
							Base: token.Token{Type: token.IDENT, Text: "int"},
						},
					},
				},
			},
			{
				input: "type Array [2]int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Array"},
						DataType: &ast.ArrayDataType{
							LBrack:   token.Token{Type: token.LBRACK, Text: "["},
							Length:   token.Token{Type: token.INT, Text: "2"},
							RBrack:   token.Token{Type: token.RBRACK, Text: "]"},
							DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
				},
			},
			{
				input: "type Array [...]int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Array"},
						DataType: &ast.ArrayDataType{
							LBrack:   token.Token{Type: token.LBRACK, Text: "["},
							Length:   token.Token{Type: token.ELLIPSIS, Text: "..."},
							RBrack:   token.Token{Type: token.RBRACK, Text: "]"},
							DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
				},
			},
			{
				input: "type Map map[string]int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Map"},
						DataType: &ast.MapDataType{
							Map:    token.Token{Type: token.MAP, Text: "map"},
							LBrack: token.Token{Type: token.LBRACK, Text: "["},
							Key:    &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "string"}},
							RBrack: token.Token{Type: token.RBRACK, Text: "]"},
							Value:  &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
				},
			},
			{
				input: "type Pointer *int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Pointer"},
						DataType: &ast.PointerDataType{
							Star:     token.Token{Type: token.MUL, Text: "*"},
							DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
				},
			},
			{
				input: "type Slice []int",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Slice"},
						DataType: &ast.SliceDataType{
							LBrack:   token.Token{Type: token.LBRACK, Text: "["},
							RBrack:   token.Token{Type: token.RBRACK, Text: "]"},
							DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
				},
			},
			{
				input: "type Foo {}",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Foo"},
						DataType: &ast.StructDataType{
							LBrace: token.Token{Type: token.LBRACE, Text: "{"},
							RBrace: token.Token{Type: token.RBRACE, Text: "}"},
						},
					},
				},
			},
			{
				input: "type Foo {Name string}",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Foo"},
						DataType: &ast.StructDataType{
							LBrace: token.Token{Type: token.LBRACE, Text: "{"},
							Elements: ast.ElemExprList{
								{
									Name: []token.Token{{Type: token.IDENT, Text: "Name"}},
									DataType: &ast.BaseDataType{
										Base: token.Token{Type: token.IDENT, Text: "string"},
									},
								},
							},
							RBrace: token.Token{Type: token.RBRACE, Text: "}"},
						},
					},
				},
			},
			{
				input: "type Foo {Name,Desc string}",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Foo"},
						DataType: &ast.StructDataType{
							LBrace: token.Token{Type: token.LBRACE, Text: "{"},
							Elements: ast.ElemExprList{
								{
									Name: []token.Token{{Type: token.IDENT, Text: "Name"}, {Type: token.IDENT, Text: "Desc"}},
									DataType: &ast.BaseDataType{
										Base: token.Token{Type: token.IDENT, Text: "string"},
									},
								},
							},
							RBrace: token.Token{Type: token.RBRACE, Text: "}"},
						},
					},
				},
			},
			{
				input: "type Foo {Name string\n Age int `json:\"age\"`}",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Foo"},
						DataType: &ast.StructDataType{
							LBrace: token.Token{Type: token.LBRACE, Text: "{"},
							Elements: ast.ElemExprList{
								{
									Name: []token.Token{{Type: token.IDENT, Text: "Name"}},
									DataType: &ast.BaseDataType{
										Base: token.Token{Type: token.IDENT, Text: "string"},
									},
								},
								{
									Name: []token.Token{{Type: token.IDENT, Text: "Age"}},
									DataType: &ast.BaseDataType{
										Base: token.Token{Type: token.IDENT, Text: "int"},
									},
									Tag: token.Token{Type: token.RAW_STRING, Text: "`json:\"age\"`"},
								},
							},
							RBrace: token.Token{Type: token.RBRACE, Text: "}"},
						},
					},
				},
			},
			{
				input: "type Foo {Bar {Name string}}",
				expected: &ast.TypeLiteralStmt{
					Type: token.Token{Type: token.TYPE, Text: "type"},
					Expr: &ast.TypeExpr{
						Name: token.Token{Type: token.IDENT, Text: "Foo"},
						DataType: &ast.StructDataType{
							LBrace: token.Token{Type: token.LBRACE, Text: "{"},
							Elements: ast.ElemExprList{
								{
									Name: []token.Token{{Type: token.IDENT, Text: "Bar"}},
									DataType: &ast.StructDataType{
										LBrace: token.Token{Type: token.LBRACE, Text: "{"},
										Elements: ast.ElemExprList{
											{
												Name: []token.Token{{Type: token.IDENT, Text: "Name"}},
												DataType: &ast.BaseDataType{
													Base: token.Token{Type: token.IDENT, Text: "string"},
												},
											},
										},
										RBrace: token.Token{Type: token.RBRACE, Text: "}"},
									},
								},
							},
							RBrace: token.Token{Type: token.RBRACE, Text: "}"},
						},
					},
				},
			},
		}
		for _, val := range testData {
			p := New("test.api", val.input, SkipComment)
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
					Type:   token.Token{Type: token.TYPE, Text: "type"},
					LParen: token.Token{Type: token.LPAREN, Text: "("},
					ExprList: []*ast.TypeExpr{
						{
							Name:     token.Token{Type: token.IDENT, Text: "Int"},
							DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
					RParen: token.Token{Type: token.RPAREN, Text: ")"},
				},
			},
			{
				input: "type (Int = int)",
				expected: &ast.TypeGroupStmt{
					Type:   token.Token{Type: token.TYPE, Text: "type"},
					LParen: token.Token{Type: token.LPAREN, Text: "("},
					ExprList: []*ast.TypeExpr{
						{
							Name:     token.Token{Type: token.IDENT, Text: "Int"},
							Assign:   token.Token{Type: token.ASSIGN, Text: "="},
							DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
						},
					},
					RParen: token.Token{Type: token.RPAREN, Text: ")"},
				},
			},
			{
				input: "type (Array [2]int)",
				expected: &ast.TypeGroupStmt{
					Type:   token.Token{Type: token.TYPE, Text: "type"},
					LParen: token.Token{Type: token.LPAREN, Text: "("},
					ExprList: []*ast.TypeExpr{
						{
							Name: token.Token{Type: token.IDENT, Text: "Array"},
							DataType: &ast.ArrayDataType{
								LBrack:   token.Token{Type: token.LBRACK, Text: "["},
								Length:   token.Token{Type: token.INT, Text: "2"},
								RBrack:   token.Token{Type: token.RBRACK, Text: "]"},
								DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
							},
						},
					},
					RParen: token.Token{Type: token.RPAREN, Text: ")"},
				},
			},
			{
				input: "type (Array [...]int)",
				expected: &ast.TypeGroupStmt{
					Type:   token.Token{Type: token.TYPE, Text: "type"},
					LParen: token.Token{Type: token.LPAREN, Text: "("},
					ExprList: []*ast.TypeExpr{
						{
							Name: token.Token{Type: token.IDENT, Text: "Array"},
							DataType: &ast.ArrayDataType{
								LBrack:   token.Token{Type: token.LBRACK, Text: "["},
								Length:   token.Token{Type: token.ELLIPSIS, Text: "..."},
								RBrack:   token.Token{Type: token.RBRACK, Text: "]"},
								DataType: &ast.BaseDataType{Base: token.Token{Type: token.IDENT, Text: "int"}},
							},
						},
					},
					RParen: token.Token{Type: token.RPAREN, Text: ")"},
				},
			},
		}
		for _, val := range testData {
			p := New("test.api", val.input, SkipComment)
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
	}

	for _, v := range testData {
		p := New("test.api", v, SkipComment)
		p.Parse()
		assertx.ErrorOrigin(t, v, p.errors...)
	}
}
