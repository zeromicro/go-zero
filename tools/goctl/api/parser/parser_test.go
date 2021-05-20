package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testApi = "// syntax doc\nsyntax = \"v1\" // syntax comment\n\n// type doc\ntype Request {\n\tName string `path:\"name,options=you|me\"`\n}\n\ntype Response {\n\tMessage string `json:\"message\"`\n}\n\n// service doc\nservice greet-api {\n\t// handler doc\n\t@handler GreetHandler // handler comment\n\tget /from/:name(Request) returns (Response);\n}"

func TestParseContent(t *testing.T) {
	spec, err := ParseContent(testApi)
	assert.Nil(t, err)
	assert.Equal(t, []string{`// syntax doc`}, spec.Syntax.Doc)
	assert.Equal(t, []string{`// syntax comment`}, spec.Syntax.Comment)
	for _, tp := range spec.Types {
		if tp.Name() == "Request" {
			assert.Equal(t, []string{`// type doc`}, tp.Documents())
		}
	}
	for _, e := range spec.Service.Routes() {
		if e.Handler == "GreetHandler" {
			assert.Equal(t, []string{"// handler doc"}, e.HandlerDoc)
			assert.Equal(t, []string{"// handler comment"}, e.HandlerComment)
		}
	}
}
