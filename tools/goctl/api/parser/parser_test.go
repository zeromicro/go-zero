package parser

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

//go:embed testdata/test.api
var testApi string

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

func TestMissingService(t *testing.T) {
	sp, err := ParseContent("")
	assert.Nil(t, err)
	err = sp.Validate()
	assert.Equal(t, spec.ErrMissingService, err)
}
