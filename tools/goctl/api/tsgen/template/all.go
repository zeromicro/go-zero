package template

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed components.tpl
var Components string

//go:embed nested.tpl
var Nested string

//go:embed handlers.tpl
var Handlers string

//go:embed requests.tpl
var Requests string

type ComponentMemberTemplateData struct {
	OptionalTag  string
	PropertyName string
	PropertyType string
	Docs         []string
	Comment      string
}

type ComponentTypeTemplateData struct {
	TypeName string
	Members  []*ComponentMemberTemplateData
	SubTypes []*ComponentTypeTemplateData
}

type ComponentTemplateData struct {
	Version string
	Types   []ComponentTypeTemplateData
}

type ComponentNestedTypeTemplateData struct {
	TypeName string
	Indent   int
	Members  []*ComponentMemberTemplateData
}

type HandlerTemplateData struct {
	Caller        string
	IsUnwrapAPI   bool
	ComponentName string
	Routes        []*HandlerRouteTemplateData
}

type HandlerRouteTemplateData struct {
	Comment       string
	HttpMethod    string
	FuncName      string
	FuncArgs      string
	GenericsTypes string
	ResponseType  string
	CallArgs      string
}

type RequestTemplateData struct {
	Caller string
}

func indent(n int) string {
	return strings.Repeat(" ", n)
}

func GenTs(writer io.Writer, tpl string, data any) error {
	tmp, err := template.New("tmp").
		Funcs(template.FuncMap{
			"Indent": indent,
		}).
		Parse(tpl)
	if err != nil {
		return err
	}

	return tmp.Execute(writer, data)
}

func GenTsFile(dir string, name string, tpl string, data any) error {
	tmp, err := template.New(name).
		Funcs(template.FuncMap{
			"Indent": indent,
		}).
		Parse(tpl)
	if err != nil {
		return err
	}
	p := filepath.Join(dir, fmt.Sprintf("%s.ts", name))
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmp.Execute(f, data)
}
