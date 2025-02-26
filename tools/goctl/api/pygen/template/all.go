package template

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/util"
)

//go:embed base.tpl
var ApiBase string

//go:embed client.tpl
var ApiClient string

//go:embed message.tpl
var ApiMessage string

type PyFieldTemplateData struct {
	FieldName string
	// FieldType    string
	FieldTag     string
	FieldTagName string
}

type PyMessageTemplateData struct {
	MessageName string
	Fields      []*PyFieldTemplateData
	HeaderCount int
	BodyCount   int
	PathCount   int
	FormCount   int
}

type PyMessagesTemplateData struct {
	Messages []*PyMessageTemplateData
}

type PyActionTemplateData struct {
	ActionName      string
	HttpMethod      string
	UrlPrefix       string
	UrlPath         string
	RequestMessage  *PyMessageTemplateData
	ResponseMessage *PyMessageTemplateData
}

type PyClientTemplateData struct {
	ClientName string
	Actions    []*PyActionTemplateData
}

type PyBaseTemplateData struct {
	ClientName string
}

func indent(n int) string {
	return strings.Repeat(" ", n)
}

func GenFile(dir string, filename string, tpl string, data any) error {
	tmp, err := template.New(filename).
		Funcs(template.FuncMap{
			"indent":    indent,
			"SnakeCase": util.SnakeCase,
			"ToLower":   strings.ToLower,
			"ToUpper":   strings.ToUpper,
		}).
		Parse(tpl)
	if err != nil {
		return err
	}

	p := filepath.Join(dir, filename)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmp.Execute(f, data)
}
