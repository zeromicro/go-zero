package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/api/parser/g4/ast"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const infoBlock = `
info(
	title: "foo"
	desc: "bar"
)
`

func TestInfo(t *testing.T) {
	testInfo(t, map[string]string{
		"title": "foo",
		"desc":  "bar",
	}, false, `info(
		title: "foo"
		desc: "bar"
	)`)

	testInfo(t, map[string]string{}, false, `info()`)
	testInfo(t, map[string]string{
		"title": "",
		"desc":  "",
	}, false, `info(
		title:
		desc:
	)`)

	testInfo(t, map[string]string{
		"title": "foo",
		"desc":  "foo\n\t\tbar",
	}, false, `info(
		title: "foo"
		desc: "foo
		bar"		
	)`)

	testInfo(t, nil, true, `info`)
	testInfo(t, nil, true, `info (`)
}

func testInfo(t *testing.T, expected interface{}, expectErr bool, content string) {
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
	result := p.InfoBlock().Accept(visitor)

	imp, ok := result.(*spec.Info)
	assert.True(t, ok)
	assert.Equal(t, expected, imp.Proterties)
}
