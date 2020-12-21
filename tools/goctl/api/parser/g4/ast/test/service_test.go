package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

func TestService(t *testing.T) {
	testServiceAnnotation(t, "", true, `
	@server(
		jwt: Foo
		group: foo/bar
		anotherKey: anotherValue
	)
	service example-api {
	}
	`, "jwt", "Foo")
}

func testServiceAnnotation(t *testing.T, expected interface{}, expectedParserErr bool, content, key, value string) {
	var parserErr error
	p := ast.NewParser(content, ast.WithErrorCallback(func(err error) {
		if expectedParserErr {
			parserErr = err
			assert.Error(t, err)
			if logEnable {
				fmt.Printf("%+v\r\n", err)
			}
			return
		}
		assert.Nil(t, err)
	}))
	visitor := ast.NewApiVisitor()
	result := p.Annotation().Accept(visitor)
	if parserErr == nil {
		anno, ok := result.(*spec.Annotation)
		assert.True(t, ok)
		assert.Equal(t, anno.Properties[key], value)
	}
}
