package util

import (
	"testing"

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
