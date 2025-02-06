package template

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/csgen/util"
)

//go:embed ApiAttribute.tpl
var ApiAttribute string

//go:embed ApiBodyJsonContent.tpl
var ApiBodyJsonContent string

//go:embed ApiException.tpl
var ApiException string

//go:embed ApiBaseClient.tpl
var ApiBaseClient string

//go:embed ApiClient.tpl
var ApiClient string

//go:embed ApiMessage.tpl
var ApiMessage string

type CSharpTemplateData struct {
	Namespace string
}

type CSharpApiMessageFieldTemplateData struct {
	FieldName  string
	KeyName    string
	TypeName   string
	Tag        string
	IsOptional bool
}

type CSharpApiMessageTemplateData struct {
	CSharpTemplateData
	MessageName string
	Fields      []CSharpApiMessageFieldTemplateData
}

type CSharpApiClientRouteTemplateData struct {
	HttpMethod   string
	Prefix       string
	ActionPrefix string
	ActionName   string
	UrlPath      string
	RequestType  *string
	ResponseType *string
}

type CSharpApiClientTemplateData struct {
	CSharpTemplateData
	ClientName string
	Routes     []CSharpApiClientRouteTemplateData
}

func WriteFile(dir string, name string, tpl string, data any) error {
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"CamelCase": util.CamelCase,
			"ToUpper":   strings.ToUpper,
		}).
		Parse(tpl)
	if err != nil {
		return err
	}
	p := filepath.Join(dir, fmt.Sprintf("%s.cs", name))
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
