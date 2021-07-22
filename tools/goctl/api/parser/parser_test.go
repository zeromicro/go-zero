package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

var testApi = "// syntax doc\nsyntax = \"v1\" // syntax comment\n\n// type doc\ntype Request {\n\tName string `path:\"name,options=you|me\"`\n}\n\ntype Response {\n\tMessage string `json:\"message\"`\n}\n\n// service doc\nservice greet-api {\n\t// handler doc\n\t@handler GreetHandler // handler comment\n\tget /from/:name(Request) returns (Response);\n}"

func TestParseContent(t *testing.T) {
	sp, err := ParseContent(testApi)
	assert.Nil(t, err)
	assert.Equal(t, spec.Doc{`// syntax doc`}, sp.Syntax.Doc)
	assert.Equal(t, spec.Doc{`// syntax comment`}, sp.Syntax.Comment)
	for _, tp := range sp.Types {
		if tp.Name() == "Request" {
			assert.Equal(t, []string{`// type doc`}, tp.Documents())
		}
	}
	for _, e := range sp.Service.Routes() {
		if e.Handler == "GreetHandler" {
			assert.Equal(t, spec.Doc{"// handler doc"}, e.HandlerDoc)
			assert.Equal(t, spec.Doc{"// handler comment"}, e.HandlerComment)
		}
	}
}
