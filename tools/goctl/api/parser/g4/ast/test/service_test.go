package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

func TestService(t *testing.T) {
	testServiceAnnotation(t, `@server(
		jwt: Foo
		group: foo/bar
		anotherKey: anotherValue
	)
	`, "jwt", "Foo")
}

func testServiceAnnotation(t *testing.T, content, key, value string) {
	var parserErr error
	p := ast.NewParser(content, ast.WithErrorCallback(func(err error) {
		assert.Nil(t, err)
	}))
	visitor := ast.NewApiVisitor()
	result := p.ServerMeta().Accept(visitor)
	if parserErr == nil {
		anno, ok := result.(spec.Annotation)
		assert.True(t, ok)
		assert.Equal(t, anno.Properties[key], value)
	}
}
