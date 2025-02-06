package template

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/unigen/util"
)

//go:embed ApiBaseClient.tpl
var ApiBaseClient string

//go:embed ApiClient.tpl
var ApiClient string

//go:embed ApiMessage.tpl
var ApiMessage string

type UniAppApiClientRouteTemplateData struct {
	HttpMethod            string
	Prefix                string
	ActionPrefix          string
	ActionName            string
	UrlPath               string
	RequestType           *string
	RequestHasQueryString bool
	RequestHasHeaders     bool
	RequestHasBody        bool
	ResponseType          *string
	ResponseHeadersType   *string
	ResponseBodyType      *string
}

type UniAppApiClientTemplateData struct {
	ClientName       string
	RequestTypes     []string
	ResponseTypes    []string
	ResponseSubTypes map[string][]string
	Routes           []UniAppApiClientRouteTemplateData
}

type UniAppApiMessageFieldTemplateData struct {
	FieldName  string
	TypeName   string
	IsOptional bool
}

type UniAppApiSubMessageTemplateData struct {
	MessageName string
	Fields      []UniAppApiMessageFieldTemplateData
}

type UniAppApiMessageTemplateData struct {
	MessageName string
	Fields      []UniAppApiMessageFieldTemplateData
	SubMessages []UniAppApiSubMessageTemplateData
	ImportTypes []string
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
	p := filepath.Join(dir, fmt.Sprintf("%s.ts", name))
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
