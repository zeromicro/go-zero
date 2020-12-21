package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const importLit = `import "foo.api"`

const importGroup = `
import (
	"foo.api"
	"bar.api"
	"foo/bar.api"
)
`

func TestImport(t *testing.T) {
	testImport(t, []string{"foo.api"}, false, importLit)
	testImport(t, []string{"foo.api", "bar.api", "foo/bar.api"}, false, importGroup)
	testImport(t, nil, false, `import ()`)
	testImport(t, nil, true, `import`)
	testImport(t, nil, true, `import user.api`)
	testImport(t, nil, true, `import "user.api "`)
	testImport(t, nil, true, `import "/"`)
	testImport(t, nil, true, `import " "`)
	testImport(t, nil, true, `import "user-.api"`)
}

func testImport(t *testing.T, expected []string, expectErr bool, content string) {
	defer func() {
		p := recover()
		if expectErr {
			assert.NotNil(t, p)
			return
		}
		assert.Nil(t, p)
	}()

	p := ast.NewParser(content)
	visitor := ast.NewApiVisitor()
	result := p.ImportSpec().Accept(visitor)

	imp, ok := result.(*spec.ApiImport)
	assert.True(t, ok)
	sort.Strings(imp.List)
	sort.Strings(expected)
	assert.Equal(t, expected, imp.List)
}
