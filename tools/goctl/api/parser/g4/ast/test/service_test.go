package test

import (
	"testing"
)

func TestService(t *testing.T) {
	testService(t, "", true, `
	@server(
		jwt: Foo
		group: foo/bar
		anotherKey: anotherValue
	)
	service example-api {
	}
	`)
}

func testService(t *testing.T, expected interface{}, expectedParserErr bool, content string) {
	//var parserErr error
	//p := ast.NewParser(content, ast.WithErrorCallback(func(err error) {
	//	if expectedParserErr {
	//		parserErr = err
	//		assert.Error(t, err)
	//		if logEnable {
	//			fmt.Printf("%+v\r\n", err)
	//		}
	//		return
	//	}
	//	assert.Nil(t, err)
	//}))
	//visitor := ast.NewApiVisitor()
	//result := p.ServiceBlock().Accept(visitor)
	//if parserErr == nil {
	//	visitResult, ok := result.(*ast.VisitResult)
	//	assert.True(t, ok)
	//
	//	r, err := visitResult.Result()
	//	assert.Nil(t, err)
	//	syntax, ok := r.(*spec.ApiSyntax)
	//	assert.True(t, ok)
	//	assert.Equal(t, expected, syntax.Version)
	//}
}
