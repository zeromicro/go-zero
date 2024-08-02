package test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
)

//go:embed apis/test.api
var testApi string

var parser = ast.NewParser(ast.WithParserPrefix("test.api"), ast.WithParserDebug())

func TestApi(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.Api().Accept(visitor)
	}

	v, err := parser.Accept(fn, testApi)
	assert.Nil(t, err)
	api := v.(*ast.Api)
	body := &ast.Body{
		Lp:   ast.NewTextExpr("("),
		Rp:   ast.NewTextExpr(")"),
		Name: &ast.Literal{Literal: ast.NewTextExpr("FooBar")},
	}

	returns := ast.NewTextExpr("returns")
	assert.True(t, api.Equal(&ast.Api{
		Syntax: &ast.SyntaxExpr{
			Syntax:  ast.NewTextExpr("syntax"),
			Assign:  ast.NewTextExpr("="),
			Version: ast.NewTextExpr(`"v1"`),
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// syntax doc"),
			},
			CommentExpr: ast.NewTextExpr("// syntax comment"),
		},
		Import: []*ast.ImportExpr{
			{
				Import: ast.NewTextExpr("import"),
				Value:  ast.NewTextExpr(`"foo.api"`),
				DocExpr: []ast.Expr{
					ast.NewTextExpr("// import doc"),
				},
				CommentExpr: ast.NewTextExpr("// import comment"),
			},
			{
				Import: ast.NewTextExpr("import"),
				Value:  ast.NewTextExpr(`"bar.api"`),
				DocExpr: []ast.Expr{
					ast.NewTextExpr("// import group doc"),
				},
				CommentExpr: ast.NewTextExpr("// import group comment"),
			},
		},
		Info: &ast.InfoExpr{
			Info: ast.NewTextExpr("info"),
			Lp:   ast.NewTextExpr("("),
			Rp:   ast.NewTextExpr(")"),
			Kvs: []*ast.KvExpr{
				{
					Key:   ast.NewTextExpr("author"),
					Value: ast.NewTextExpr(`"songmeizi"`),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// author doc"),
					},
					CommentExpr: ast.NewTextExpr("// author comment"),
				},
				{
					Key:   ast.NewTextExpr("date"),
					Value: ast.NewTextExpr(`2020-01-04`),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// date doc"),
					},
					CommentExpr: ast.NewTextExpr("// date comment"),
				},
				{
					Key: ast.NewTextExpr("desc"),
					Value: ast.NewTextExpr(`"break line
    desc"`),
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// desc doc"),
					},
					CommentExpr: ast.NewTextExpr("// desc comment"),
				},
			},
		},
		Type: []ast.TypeExpr{
			&ast.TypeStruct{
				Name:   ast.NewTextExpr("FooBar"),
				Struct: ast.NewTextExpr("struct"),
				LBrace: ast.NewTextExpr("{"),
				RBrace: ast.NewTextExpr("}"),
				Fields: []*ast.TypeField{
					{
						Name:     ast.NewTextExpr("Foo"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
					},
				},
			},
			&ast.TypeStruct{
				Name:   ast.NewTextExpr("Bar"),
				LBrace: ast.NewTextExpr("{"),
				RBrace: ast.NewTextExpr("}"),
				DocExpr: []ast.Expr{
					ast.NewTextExpr("// remove struct"),
				},
				Fields: []*ast.TypeField{
					{
						Name:     ast.NewTextExpr("VString"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("string")},
						Tag:      ast.NewTextExpr("`json:\"vString\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vString"),
						},
					},
					{
						Name:     ast.NewTextExpr("VBool"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("bool")},
						Tag:      ast.NewTextExpr("`json:\"vBool\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vBool"),
						},
					},
					{
						Name:     ast.NewTextExpr("VInt8"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("int8")},
						Tag:      ast.NewTextExpr("`json:\"vInt8\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vInt8"),
						},
					},
					{
						Name:     ast.NewTextExpr("VInt16"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("int16")},
						Tag:      ast.NewTextExpr("`json:\"vInt16\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vInt16"),
						},
					},
					{
						Name:     ast.NewTextExpr("VInt32"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("int32")},
						Tag:      ast.NewTextExpr("`json:\"vInt32\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vInt32"),
						},
					},
					{
						Name:     ast.NewTextExpr("VInt64"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("int64")},
						Tag:      ast.NewTextExpr("`json:\"vInt64\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vInt64"),
						},
					},
					{
						Name:     ast.NewTextExpr("VInt"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
						Tag:      ast.NewTextExpr("`json:\"vInt\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vInt"),
						},
					},
					{
						Name:     ast.NewTextExpr("VUInt8"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("uint8")},
						Tag:      ast.NewTextExpr("`json:\"vUInt8\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vUInt8"),
						},
					},
					{
						Name:     ast.NewTextExpr("VUInt16"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("uint16")},
						Tag:      ast.NewTextExpr("`json:\"vUInt16\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vUInt16"),
						},
					},
					{
						Name:     ast.NewTextExpr("VUInt32"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("uint32")},
						Tag:      ast.NewTextExpr("`json:\"vUInt32\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vUInt32"),
						},
					},
					{
						Name:     ast.NewTextExpr("VUInt64"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("uint64")},
						Tag:      ast.NewTextExpr("`json:\"vUInt64\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vUInt64"),
						},
					},
					{
						Name:     ast.NewTextExpr("VFloat32"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("float32")},
						Tag:      ast.NewTextExpr("`json:\"vFloat32\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vFloat32"),
						},
					},
					{
						Name:     ast.NewTextExpr("VFloat64"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("float64")},
						Tag:      ast.NewTextExpr("`json:\"vFloat64\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vFloat64"),
						},
					},
					{
						Name:     ast.NewTextExpr("VByte"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("byte")},
						Tag:      ast.NewTextExpr("`json:\"vByte\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vByte"),
						},
					},
					{
						Name:     ast.NewTextExpr("VRune"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("rune")},
						Tag:      ast.NewTextExpr("`json:\"vRune\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vRune"),
						},
					},
					{
						Name: ast.NewTextExpr("VMap"),
						DataType: &ast.Map{
							MapExpr: ast.NewTextExpr("map[string]int"),
							Map:     ast.NewTextExpr("map"),
							LBrack:  ast.NewTextExpr("["),
							RBrack:  ast.NewTextExpr("]"),
							Key:     ast.NewTextExpr("string"),
							Value:   &ast.Literal{Literal: ast.NewTextExpr("int")},
						},
						Tag: ast.NewTextExpr("`json:\"vMap\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vMap"),
						},
					},
					{
						Name: ast.NewTextExpr("VArray"),
						DataType: &ast.Array{
							ArrayExpr: ast.NewTextExpr("[]int"),
							LBrack:    ast.NewTextExpr("["),
							RBrack:    ast.NewTextExpr("]"),
							Literal:   &ast.Literal{Literal: ast.NewTextExpr("int")},
						},
						Tag: ast.NewTextExpr("`json:\"vArray\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vArray"),
						},
					},
					{
						Name:     ast.NewTextExpr("VStruct"),
						DataType: &ast.Literal{Literal: ast.NewTextExpr("FooBar")},
						Tag:      ast.NewTextExpr("`json:\"vStruct\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vStruct"),
						},
					},
					{
						Name: ast.NewTextExpr("VStructPointer"),
						DataType: &ast.Pointer{
							PointerExpr: ast.NewTextExpr("*FooBar"),
							Star:        ast.NewTextExpr("*"),
							Name:        ast.NewTextExpr("FooBar"),
						},
						Tag: ast.NewTextExpr("`json:\"vStructPointer\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vStructPointer"),
						},
					},
					{
						Name:     ast.NewTextExpr("VInterface"),
						DataType: &ast.Interface{Literal: ast.NewTextExpr("interface{}")},
						Tag:      ast.NewTextExpr("`json:\"vInterface\"`"),
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// vInterface"),
						},
					},
					{
						IsAnonymous: true,
						DataType:    &ast.Literal{Literal: ast.NewTextExpr("FooBar")},
						DocExpr: []ast.Expr{
							ast.NewTextExpr("// inline"),
						},
					},
				},
			},
		},
		Service: []*ast.Service{
			{
				AtServer: &ast.AtServer{
					AtServerToken: ast.NewTextExpr("@server"),
					Lp:            ast.NewTextExpr("("),
					Rp:            ast.NewTextExpr(")"),
					Kv: []*ast.KvExpr{
						{
							Key:   ast.NewTextExpr("host"),
							Value: ast.NewTextExpr("0.0.0.0"),
						},
						{
							Key:   ast.NewTextExpr("port"),
							Value: ast.NewTextExpr("8080"),
						},
						{
							Key: ast.NewTextExpr("annotation"),
							Value: ast.NewTextExpr(`"break line
    desc"`),
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
								Name:           ast.NewTextExpr("postFoo"),
							},
							Route: &ast.Route{
								Method:      ast.NewTextExpr("post"),
								Path:        ast.NewTextExpr("/foo"),
								Req:         body,
								ReturnToken: returns,
								Reply:       body,
								DocExpr: []ast.Expr{
									ast.NewTextExpr("// foo"),
								},
							},
						},
						{
							AtDoc: &ast.AtDoc{
								AtDocToken: ast.NewTextExpr("@doc"),
								Lp:         ast.NewTextExpr("("),
								Rp:         ast.NewTextExpr(")"),
								Kv: []*ast.KvExpr{
									{
										Key:   ast.NewTextExpr("summary"),
										Value: ast.NewTextExpr("bar"),
									},
								},
							},
							AtServer: &ast.AtServer{
								AtServerToken: ast.NewTextExpr("@server"),
								Lp:            ast.NewTextExpr("("),
								Rp:            ast.NewTextExpr(")"),
								Kv: []*ast.KvExpr{
									{
										Key:   ast.NewTextExpr("handler"),
										Value: ast.NewTextExpr("postBar"),
									},
								},
							},
							Route: &ast.Route{
								Method: ast.NewTextExpr("post"),
								Path:   ast.NewTextExpr("/bar"),
								Req:    body,
							},
						},
						{
							AtDoc: &ast.AtDoc{
								AtDocToken: ast.NewTextExpr("@doc"),
								Lp:         ast.NewTextExpr("("),
								Rp:         ast.NewTextExpr(")"),
								LineDoc:    ast.NewTextExpr(`"foobar"`),
							},
							AtHandler: &ast.AtHandler{
								AtHandlerToken: ast.NewTextExpr("@handler"),
								Name:           ast.NewTextExpr("postFooBar"),
							},
							Route: &ast.Route{
								Method:      ast.NewTextExpr("post"),
								Path:        ast.NewTextExpr("/foo/bar"),
								ReturnToken: returns,
								Reply:       body,
								DocExpr: []ast.Expr{
									ast.NewTextExpr(`/**
    * httpmethod: post
    * path: /foo/bar
    * reply: FooBar
    */`),
								},
							},
						},
						{
							AtDoc: &ast.AtDoc{
								AtDocToken: ast.NewTextExpr("@doc"),
								Lp:         ast.NewTextExpr("("),
								Rp:         ast.NewTextExpr(")"),
								LineDoc:    ast.NewTextExpr(`"barfoo"`),
							},
							AtHandler: &ast.AtHandler{
								AtHandlerToken: ast.NewTextExpr("@handler"),
								Name:           ast.NewTextExpr("postBarFoo"),
							},
							Route: &ast.Route{
								Method:      ast.NewTextExpr("post"),
								Path:        ast.NewTextExpr("/bar/foo"),
								ReturnToken: returns,
								Reply:       body,
								CommentExpr: ast.NewTextExpr("// post:/bar/foo"),
							},
						},
						{
							AtDoc: &ast.AtDoc{
								AtDocToken: ast.NewTextExpr("@doc"),
								Lp:         ast.NewTextExpr("("),
								Rp:         ast.NewTextExpr(")"),
								LineDoc:    ast.NewTextExpr(`"barfoo"`),
							},
							AtHandler: &ast.AtHandler{
								AtHandlerToken: ast.NewTextExpr("@handler"),
								Name:           ast.NewTextExpr("getBarFoo"),
							},
							Route: &ast.Route{
								Method:      ast.NewTextExpr("get"),
								Path:        ast.NewTextExpr("/bar/foo"),
								ReturnToken: returns,
								Reply:       body,
							},
						},
					},
				},
			},
		},
	}))
}

func TestApiSyntax(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.Api().Accept(visitor)
	}
	parser.Accept(fn, `
	// doc 1
	// doc 2
	syntax = "v1" // comment 1
	// comment 2
	import "foo.api"
	`)
}
