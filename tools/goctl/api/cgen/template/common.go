package template

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/util"
)

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

	if _, err := tmp.New("to_primitive").Parse(cJsonToPrimitiveTemplate); err != nil {
		return err
	}

	if _, err := tmp.New("to_array").Parse(cJsonToArrayTemplate); err != nil {
		return err
	}

	if _, err := tmp.New("to_object").Parse(cJsonToObjectTemplate); err != nil {
		return err
	}

	if _, err := tmp.New("from_primitive").Parse(cJsonFromPrimitiveTemplate); err != nil {
		return err
	}

	if _, err := tmp.New("from_array").Parse(cJsonFromArrayTemplate); err != nil {
		return err
	}

	if _, err := tmp.New("from_object").Parse(cJsonFromObjectTemplate); err != nil {
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
