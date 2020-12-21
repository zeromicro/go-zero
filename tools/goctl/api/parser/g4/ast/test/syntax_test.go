package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
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

func testSyntax(t *testing.T, expected interface{}, expectedParserErr bool, content string) {
	var globalErr error
	p := ast.NewParser(content, ast.WithErrorCallback(func(err error) {
		globalErr = err
		if expectedParserErr {
			assert.Error(t, err)
			fmt.Printf("%+v\r\n", err)
			return
		}
		assert.Nil(t, err)
	}))
	visitor := ast.NewApiVisitor()
	result := p.SyntaxLit().Accept(visitor)
	if globalErr != nil {
		return
	}

	visitResult, ok := result.(*ast.VisitResult)
	assert.True(t, ok)
	r, err := visitResult.Result()
	assert.Nil(t, err)

	syntax, ok := r.(*ast.Syntax)
	assert.True(t, ok)
	assert.Equal(t, expected, syntax.Version)
}
