package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
)

var fieldAccept = func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
	return p.Field().Accept(visitor)
}

func TestField(t *testing.T) {
	t.Run("anonymous", func(t *testing.T) {
		v, err := parser.Accept(fieldAccept, `User`)
		assert.Nil(t, err)
		f := v.(*ast.TypeField)
		assert.True(t, f.Equal(&ast.TypeField{
			IsAnonymous: true,
			DataType:    &ast.Literal{Literal: ast.NewTextExpr("User")},
		}))

		v, err = parser.Accept(fieldAccept, `*User`)
		assert.Nil(t, err)
		f = v.(*ast.TypeField)
		assert.True(t, f.Equal(&ast.TypeField{
			IsAnonymous: true,
			DataType: &ast.Pointer{
				PointerExpr: ast.NewTextExpr("*User"),
				Star:        ast.NewTextExpr("*"),
				Name:        ast.NewTextExpr("User"),
			},
		}))

		v, err = parser.Accept(fieldAccept, `
		// anonymous user
		*User // pointer type`)
		assert.Nil(t, err)
		f = v.(*ast.TypeField)
		assert.True(t, f.Equal(&ast.TypeField{
			IsAnonymous: true,
			DataType: &ast.Pointer{
				PointerExpr: ast.NewTextExpr("*User"),
				Star:        ast.NewTextExpr("*"),
				Name:        ast.NewTextExpr("User"),
			},
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// anonymous user"),
			},
			CommentExpr: ast.NewTextExpr("// pointer type"),
		}))

		_, err = parser.Accept(fieldAccept, `interface`)
		assert.Error(t, err)

		_, err = parser.Accept(fieldAccept, `map`)
		assert.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fieldAccept, `User int`)
		assert.Nil(t, err)
		f := v.(*ast.TypeField)
		assert.True(t, f.Equal(&ast.TypeField{
			Name:     ast.NewTextExpr("User"),
			DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
		}))
		v, err = parser.Accept(fieldAccept, `Foo Bar`)
		assert.Nil(t, err)
		f = v.(*ast.TypeField)
		assert.True(t, f.Equal(&ast.TypeField{
			Name:     ast.NewTextExpr("Foo"),
			DataType: &ast.Literal{Literal: ast.NewTextExpr("Bar")},
		}))

		v, err = parser.Accept(fieldAccept, `Foo map[int]Bar`)
		assert.Nil(t, err)
		f = v.(*ast.TypeField)
		assert.True(t, f.Equal(&ast.TypeField{
			Name: ast.NewTextExpr("Foo"),
			DataType: &ast.Map{
				MapExpr: ast.NewTextExpr("map[int]Bar"),
				Map:     ast.NewTextExpr("map"),
				LBrack:  ast.NewTextExpr("["),
				RBrack:  ast.NewTextExpr("]"),
				Key:     ast.NewTextExpr("int"),
				Value:   &ast.Literal{Literal: ast.NewTextExpr("Bar")},
			},
		}))
	})
}

func TestDataType_ID(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.DataType().Accept(visitor)
	}
	t.Run("Struct", func(t *testing.T) {
		v, err := parser.Accept(dt, `Foo`)
		assert.Nil(t, err)
		id := v.(ast.DataType)
		assert.True(t, id.Equal(&ast.Literal{Literal: ast.NewTextExpr("Foo")}))
	})

	t.Run("basic", func(t *testing.T) {
		v, err := parser.Accept(dt, `int`)
		assert.Nil(t, err)
		id := v.(ast.DataType)
		assert.True(t, id.Equal(&ast.Literal{Literal: ast.NewTextExpr("int")}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `map`)
		assert.Error(t, err)
	})
}

func TestDataType_Map(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.MapType().Accept(visitor)
	}
	t.Run("basicKey", func(t *testing.T) {
		v, err := parser.Accept(dt, `map[int]Bar`)
		assert.Nil(t, err)
		m := v.(ast.DataType)
		assert.True(t, m.Equal(&ast.Map{
			MapExpr: ast.NewTextExpr("map[int]Bar"),
			Map:     ast.NewTextExpr("map"),
			LBrack:  ast.NewTextExpr("["),
			RBrack:  ast.NewTextExpr("]"),
			Key:     ast.NewTextExpr("int"),
			Value:   &ast.Literal{Literal: ast.NewTextExpr("Bar")},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `map[var]Bar`)
		assert.Error(t, err)

		_, err = parser.Accept(dt, `map[*User]Bar`)
		assert.Error(t, err)

		_, err = parser.Accept(dt, `map[User]Bar`)
		assert.Error(t, err)
	})
}

func TestDataType_Array(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.ArrayType().Accept(visitor)
	}
	t.Run("basic", func(t *testing.T) {
		v, err := parser.Accept(dt, `[]int`)
		assert.Nil(t, err)
		array := v.(ast.DataType)
		assert.True(t, array.Equal(&ast.Array{
			ArrayExpr: ast.NewTextExpr("[]int"),
			LBrack:    ast.NewTextExpr("["),
			RBrack:    ast.NewTextExpr("]"),
			Literal:   &ast.Literal{Literal: ast.NewTextExpr("int")},
		}))
	})

	t.Run("pointer", func(t *testing.T) {
		v, err := parser.Accept(dt, `[]*User`)
		assert.Nil(t, err)
		array := v.(ast.DataType)
		assert.True(t, array.Equal(&ast.Array{
			ArrayExpr: ast.NewTextExpr("[]*User"),
			LBrack:    ast.NewTextExpr("["),
			RBrack:    ast.NewTextExpr("]"),
			Literal: &ast.Pointer{
				PointerExpr: ast.NewTextExpr("*User"),
				Star:        ast.NewTextExpr("*"),
				Name:        ast.NewTextExpr("User"),
			},
		}))
	})

	t.Run("interface{}", func(t *testing.T) {
		v, err := parser.Accept(dt, `[]interface{}`)
		assert.Nil(t, err)
		array := v.(ast.DataType)
		assert.True(t, array.Equal(&ast.Array{
			ArrayExpr: ast.NewTextExpr("[]interface{}"),
			LBrack:    ast.NewTextExpr("["),
			RBrack:    ast.NewTextExpr("]"),
			Literal:   &ast.Interface{Literal: ast.NewTextExpr("interface{}")},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `[]var`)
		assert.Error(t, err)

		_, err = parser.Accept(dt, `[]interface`)
		assert.Error(t, err)
	})
}

func TestDataType_Interface(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.DataType().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(dt, `interface{}`)
		assert.Nil(t, err)
		inter := v.(ast.DataType)
		assert.True(t, inter.Equal(&ast.Interface{Literal: ast.NewTextExpr("interface{}")}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `interface`)
		assert.Error(t, err)
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `interface{`)
		assert.Error(t, err)
	})
}

func TestDataType_Time(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.DataType().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(dt, `time.Time`)
		assert.Error(t, err)
	})
}

func TestDataType_Pointer(t *testing.T) {
	dt := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.PointerType().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(dt, `*int`)
		assert.Nil(t, err)
		assert.True(t, v.(ast.DataType).Equal(&ast.Pointer{
			PointerExpr: ast.NewTextExpr("*int"),
			Star:        ast.NewTextExpr("*"),
			Name:        ast.NewTextExpr("int"),
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(dt, `int`)
		assert.Error(t, err)
	})
}

func TestAlias(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.TypeAlias().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `Foo int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `Foo=int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `
		Foo int // comment`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `
		Foo int /**comment*/`)
		assert.Error(t, err)
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `Foo var`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `Foo 2`)
		assert.Error(t, err)
	})
}

func TestTypeStruct(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.TypeStruct().Accept(visitor)
	}

	t.Run("normal", func(t *testing.T) {
		v, err := parser.Accept(fn, "Foo {\n\t\t\tFoo string\n\t\t\tBar int `json:\"bar\"``\n\t\t}")
		assert.Nil(t, err)
		s := v.(*ast.TypeStruct)
		assert.True(t, s.Equal(&ast.TypeStruct{
			Name:   ast.NewTextExpr("Foo"),
			LBrace: ast.NewTextExpr("{"),
			RBrace: ast.NewTextExpr("}"),
			Fields: []*ast.TypeField{
				{
					Name:     ast.NewTextExpr("Foo"),
					DataType: &ast.Literal{Literal: ast.NewTextExpr("string")},
				},
				{
					Name:     ast.NewTextExpr("Bar"),
					DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
					Tag:      ast.NewTextExpr("`json:\"bar\"`"),
				},
			},
		}))

		v, err = parser.Accept(fn, "Foo struct{\n\t\t\tFoo string\n\t\t\tBar int `json:\"bar\"``\n\t\t}")
		assert.Nil(t, err)
		s = v.(*ast.TypeStruct)
		assert.True(t, s.Equal(&ast.TypeStruct{
			Name:   ast.NewTextExpr("Foo"),
			LBrace: ast.NewTextExpr("{"),
			RBrace: ast.NewTextExpr("}"),
			Struct: ast.NewTextExpr("struct"),
			Fields: []*ast.TypeField{
				{
					Name:     ast.NewTextExpr("Foo"),
					DataType: &ast.Literal{Literal: ast.NewTextExpr("string")},
				},
				{
					Name:     ast.NewTextExpr("Bar"),
					DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
					Tag:      ast.NewTextExpr("`json:\"bar\"`"),
				},
			},
		}))
	})
}

func TestTypeBlock(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.TypeBlock().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `type(
			// doc
			Foo int
		)`)
		assert.Error(t, err)

		v, err := parser.Accept(fn, `type (
			// doc
			Foo {
				Bar int
			}
		)`)
		assert.Nil(t, err)
		st := v.([]ast.TypeExpr)
		assert.True(t, st[0].Equal(&ast.TypeStruct{
			Name:   ast.NewTextExpr("Foo"),
			LBrace: ast.NewTextExpr("{"),
			RBrace: ast.NewTextExpr("}"),
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// doc"),
			},
			Fields: []*ast.TypeField{
				{
					Name:     ast.NewTextExpr("Bar"),
					DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
				},
			},
		}))
	})
}

func TestTypeLit(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.TypeLit().Accept(visitor)
	}
	t.Run("normal", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `type Foo = int`)
		assert.Error(t, err)

		_, err = parser.Accept(fn, `
		// doc
		type Foo = int // comment`)
		assert.Error(t, err)

		v, err := parser.Accept(fn, `
		// doc
		type Foo {// comment
			Bar int
		}`)
		assert.Nil(t, err)
		st := v.(*ast.TypeStruct)
		assert.True(t, st.Equal(&ast.TypeStruct{
			Name: ast.NewTextExpr("Foo"),
			Fields: []*ast.TypeField{
				{
					Name:     ast.NewTextExpr("Bar"),
					DataType: &ast.Literal{Literal: ast.NewTextExpr("int")},
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// comment"),
					},
				},
			},
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// doc"),
			},
		}))

		v, err = parser.Accept(fn, `
		// doc
		type Foo {// comment
			Bar
		}`)
		assert.Nil(t, err)
		st = v.(*ast.TypeStruct)
		assert.True(t, st.Equal(&ast.TypeStruct{
			Name: ast.NewTextExpr("Foo"),
			Fields: []*ast.TypeField{
				{
					IsAnonymous: true,
					DataType:    &ast.Literal{Literal: ast.NewTextExpr("Bar")},
					DocExpr: []ast.Expr{
						ast.NewTextExpr("// comment"),
					},
				},
			},
			DocExpr: []ast.Expr{
				ast.NewTextExpr("// doc"),
			},
		}))
	})

	t.Run("wrong", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo`)
		assert.Error(t, err)
	})
}

func TestTypeUnExported(t *testing.T) {
	fn := func(p *api.ApiParserParser, visitor *ast.ApiVisitor) any {
		return p.TypeSpec().Accept(visitor)
	}

	t.Run("type", func(t *testing.T) {
		_, err := parser.Accept(fn, `type foo {}`)
		assert.Nil(t, err)
	})

	t.Run("field", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo {
			name int
		}`)
		assert.Nil(t, err)

		_, err = parser.Accept(fn, `type Foo {
			Name int
		}`)
		assert.Nil(t, err)
	})

	t.Run("filedDataType", func(t *testing.T) {
		_, err := parser.Accept(fn, `type Foo {
			Foo *foo
			Bar []bar
			FooBar map[int]fooBar
		}`)
		assert.Nil(t, err)
	})
}
