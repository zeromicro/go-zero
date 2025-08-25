package util

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/test"
)

func TestIsTemplate(t *testing.T) {
	executor := test.NewExecutor[string, bool]()
	executor.Add([]test.Data[string, bool]{
		{
			Name: "empty",
			Want: false,
		},
		{
			Name:  "invalid",
			Input: "{foo}",
			Want:  false,
		},
		{
			Name:  "invalid",
			Input: "{.foo}",
			Want:  false,
		},
		{
			Name:  "invalid",
			Input: "$foo",
			Want:  false,
		},
		{
			Name:  "invalid",
			Input: "{{foo}}",
			Want:  false,
		},
		{
			Name:  "invalid",
			Input: "{{.}}",
			Want:  false,
		},
		{
			Name:  "valid",
			Input: "{{.foo}}",
			Want:  true,
		},
		{
			Name:  "valid",
			Input: "{{.foo.bar}}",
			Want:  true,
		},
	}...)
	executor.Run(t, IsTemplateVariable)
}

func TestTemplateVariable(t *testing.T) {
	executor := test.NewExecutor[string, string]()
	executor.Add([]test.Data[string, string]{
		{
			Name: "empty",
		},
		{
			Name:  "invalid",
			Input: "{foo}",
		},
		{
			Name:  "invalid",
			Input: "{.foo}",
		},
		{
			Name:  "invalid",
			Input: "$foo",
		},
		{
			Name:  "invalid",
			Input: "{{foo}}",
		},
		{
			Name:  "invalid",
			Input: "{{.}}",
		},
		{
			Name:  "valid",
			Input: "{{.foo}}",
			Want:  "foo",
		},
		{
			Name:  "valid",
			Input: "{{.foo.bar}}",
			Want:  "foo.bar",
		},
	}...)
	executor.Run(t, TemplateVariable)
}

func TestTemplateFuncMap(t *testing.T) {
	cases := []struct {
		name     string
		funcMap  template.FuncMap
		template string
		variable map[string]any
		want     string
	}{
		{
			name:     "upper",
			template: "{{.foo | upper}}",
			variable: map[string]any{"foo": "bar"},
			want:     "BAR",
		},
		{
			name:     "repeat upper",
			template: "{{.foo | repeat 3 | upper}}",
			variable: map[string]any{"foo": "bar"},
			want:     "BARBARBAR",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			funcMap := DefaultFuncMap()
			if c.funcMap != nil {
				funcMap = MergeWithDefaultFuncMap(c.funcMap)
			}
			tpl := template.Must(template.New("test").Funcs(funcMap).Parse(c.template))
			var buf bytes.Buffer
			err := tpl.Execute(&buf, c.variable)
			assert.NoError(t, err)
			assert.Equal(t, c.want, buf.String())
		})
	}
}
