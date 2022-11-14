package util

import (
	"bytes"
	goformat "go/format"
	"io/ioutil"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/internal/errorx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const regularPerm = 0o666

// DefaultTemplate is a tool to provides the text/template operations
type DefaultTemplate struct {
	name  string
	text  string
	goFmt bool
}

// With returns an instance of DefaultTemplate
func With(name string) *DefaultTemplate {
	return &DefaultTemplate{
		name: name,
	}
}

// Parse accepts a source template and returns DefaultTemplate
func (t *DefaultTemplate) Parse(text string) *DefaultTemplate {
	t.text = text
	return t
}

// GoFmt sets the value to goFmt and marks the generated codes will be formatted or not
func (t *DefaultTemplate) GoFmt(format bool) *DefaultTemplate {
	t.goFmt = format
	return t
}

// SaveTo writes the codes to the target path
func (t *DefaultTemplate) SaveTo(data interface{}, path string, forceUpdate bool) error {
	if pathx.FileExists(path) && !forceUpdate {
		return nil
	}

	output, err := t.Execute(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, output.Bytes(), regularPerm)
}

// Execute returns the codes after the template executed
func (t *DefaultTemplate) Execute(data interface{}) (*bytes.Buffer, error) {
	tem, err := template.New(t.name).Parse(t.text)
	if err != nil {
		return nil, errorx.Wrap(err, "template parse error:", t.text)
	}

	buf := new(bytes.Buffer)
	if err = tem.Execute(buf, data); err != nil {
		return nil, errorx.Wrap(err, "template execute error:", t.text)
	}

	if !t.goFmt {
		return buf, nil
	}

	formatOutput, err := goformat.Source(buf.Bytes())
	if err != nil {
		return nil, errorx.Wrap(err, "go format error:", buf.String())
	}

	buf.Reset()
	buf.Write(formatOutput)
	return buf, nil
}
