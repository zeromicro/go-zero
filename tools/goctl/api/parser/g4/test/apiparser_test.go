package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
)

var (
	normalApi = `
	syntax="v1"
	
	info (
		foo: bar
	)

	type Foo {
		Bar int
	}

	@server(
		foo: bar
	)
	service foo-api{
		@doc("foo")
		@handler foo
		post /foo (Foo) returns ([]int)
	}
`
	missDeclarationApi = `
	@server(
		foo: bar
	)
	service foo-api{
		@doc("foo")
		@handler foo
		post /foo (Foo) returns (Foo)
	}
	`

	missDeclarationInArrayApi = `
	@server(
		foo: bar
	)
	service foo-api{
		@doc("foo")
		@handler foo
		post /foo returns ([]Foo)
	}
	`

	missDeclarationInArrayApi2 = `
	@server(
		foo: bar
	)
	service foo-api{
		@doc("foo")
		@handler foo
		post /foo returns ([]*Foo)
	}
	`

	nestedApiImport = `
	import "foo.api"
	`

	ambiguousSyntax = `
	syntax = "v2"
	`

	ambiguousService = `
	service bar-api{
		@handler foo
		post /foo
	}
	`
	duplicateHandler = `
	service bar-api{
		@handler foo
		post /foo
	}
	`

	duplicateRoute = `
	service bar-api{
		@handler bar
		post /foo
	}
	`

	duplicateType = `
	type Foo int
	`
)

func TestApiParser(t *testing.T) {
	t.Run("missDeclarationApi", func(t *testing.T) {
		_, err := parser.ParseContent(missDeclarationApi)
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("missDeclarationApi", func(t *testing.T) {
		_, err := parser.ParseContent(missDeclarationInArrayApi)
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("missDeclarationApi", func(t *testing.T) {
		_, err := parser.ParseContent(missDeclarationInArrayApi2)
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("nestedImport", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "foo.api")
		err := ioutil.WriteFile(file, []byte(nestedApiImport), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`import "%s"`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("duplicateImport", func(t *testing.T) {
		_, err := parser.ParseContent(`
		import "foo.api"
		import "foo.api"
		`)
		assert.Error(t, err)
	})

	t.Run("duplicateKey", func(t *testing.T) {
		_, err := parser.ParseContent(`
		info (
			foo: bar
			foo: bar
		)
		`)
		assert.Error(t, err)
	})

	t.Run("ambiguousSyntax", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "foo.api")
		err := ioutil.WriteFile(file, []byte(ambiguousSyntax), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`
		syntax = "v1"
		import "%s"`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("ambiguousSyntax", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "foo.api")
		err := ioutil.WriteFile(file, []byte(ambiguousSyntax), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`
		syntax = "v1"
		import "%s"`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("ambiguousService", func(t *testing.T) {
		file := filepath.Join(t.TempDir(), "foo.api")
		err := ioutil.WriteFile(file, []byte(ambiguousService), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`
		import "%s"
		
		service foo-api{
			@handler foo
			post /foo
		}
		`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("duplicateHandler", func(t *testing.T) {
		_, err := parser.ParseContent(`
		service foo-api{
			@handler foo
			post /foo
			
			@handler foo
			post /bar
		}
		`)
		assert.Error(t, err)

		file := filepath.Join(t.TempDir(), "foo.api")
		err = ioutil.WriteFile(file, []byte(duplicateHandler), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`
		import "%s"
		service bar-api{
			@handler foo
			post /foo
		}
		`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("duplicateRoute", func(t *testing.T) {
		_, err := parser.ParseContent(`
		service foo-api{
			@handler foo
			post /foo
			
			@handler bar
			post /foo
		}
		`)
		assert.Error(t, err)

		file := filepath.Join(t.TempDir(), "foo.api")
		err = ioutil.WriteFile(file, []byte(duplicateRoute), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`
		import "%s"
		service bar-api{
			@handler foo
			post /foo
		}
		`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("duplicateType", func(t *testing.T) {
		_, err := parser.ParseContent(`
		type Foo int
		type Foo bool
		`)
		assert.Error(t, err)

		file := filepath.Join(t.TempDir(), "foo.api")
		err = ioutil.WriteFile(file, []byte(duplicateType), os.ModePerm)
		if err != nil {
			return
		}

		_, err = parser.ParseContent(fmt.Sprintf(`
		import "%s"
		
		type Foo bool
		`, file))
		assert.Error(t, err)
		fmt.Printf("%+v\n", err)
	})

	t.Run("normal", func(t *testing.T) {
		v, err := parser.ParseContent(normalApi)
		assert.Nil(t, err)
		body := &ast.Body{
			Lp:   ast.NewTextExpr("("),
			Rp:   ast.NewTextExpr(")"),
			Name: &ast.Literal{Literal: ast.NewTextExpr("Foo")},
		}

		assert.True(t, v.Equal(&ast.Api{
			Syntax: &ast.SyntaxExpr{
				Syntax:  ast.NewTextExpr("syntax"),
				Assign:  ast.NewTextExpr("="),
				Version: ast.NewTextExpr(`"v1"`),
			},
			Info: &ast.InfoExpr{
				Info: ast.NewTextExpr("info"),
				Lp:   ast.NewTextExpr("("),
				Rp:   ast.NewTextExpr(")"),
				Kvs: []*ast.KvExpr{
					{
						Key:   ast.NewTextExpr("foo"),
						Value: ast.NewTextExpr("bar"),
					},
				},
			},
			Type: []ast.TypeExpr{
				&ast.TypeStruct{
					Name:   ast.NewTextExpr("Foo"),
					LBrace: ast.NewTextExpr("{"),
					RBrace: ast.NewTextExpr("}"),
					Fields: []*ast.TypeField{
						{
							Name:     ast.NewTextExpr("Bar"),
							DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
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
								Key:   ast.NewTextExpr("foo"),
								Value: ast.NewTextExpr("bar"),
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
								},
								Route: &ast.Route{
									Method:      ast.NewTextExpr("post"),
									Path:        ast.NewTextExpr("/foo"),
									Req:         body,
									ReturnToken: ast.NewTextExpr("returns"),
									Reply: &ast.Body{
										Lp: ast.NewTextExpr("("),
										Rp: ast.NewTextExpr(")"),
										Name: &ast.Array{
											ArrayExpr: ast.NewTextExpr("[]int"),
											LBrack:    ast.NewTextExpr("["),
											RBrack:    ast.NewTextExpr("]"),
											Literal:   &ast.Literal{Literal: ast.NewTextExpr("int")},
										},
									},
								},
							},
						},
					},
				},
			},
		}))
	})
}
