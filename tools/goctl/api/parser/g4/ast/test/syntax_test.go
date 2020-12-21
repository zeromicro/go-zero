package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

func TestSyntaxLit(t *testing.T) {
	testSyntax(t, "v1", false, `syntax = "v1"`)
	testSyntax(t, "", true, `syntax = "v 1"`)
	testSyntax(t, "", true, `syntax = "1"`)
	testSyntax(t, "", true, `syntax = ""`)
	testSyntax(t, "", true, `syntax1 = "v1"`)
	testSyntax(t, "", true, `syntax`)
	testSyntax(t, "", true, `syntax=`)
	testSyntax(t, "", true, `syntax "v1"`)
	testSyntax(t, "", true, `syntax = "v0"`)
}

func testSyntax(t *testing.T, expected interface{}, exoectErr bool, content string) {
	defer func() {
		p := recover()
		if exoectErr {
			assert.NotNil(t, p)
			return
		}
		assert.Nil(t, p)
	}()
	p := ast.NewParser(content)
	visitor := ast.NewApiVisitor()
	result := p.SyntaxLit().Accept(visitor)
	syntax, ok := result.(*spec.ApiSyntax)
	assert.True(t, ok)
	assert.Equal(t, expected, syntax.Version)
}
