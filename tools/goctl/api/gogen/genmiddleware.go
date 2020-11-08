package gogen

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/util"
)

var middlewareImplementCode = `
package middleware

import "net/http"

type {{.name}} struct {
}

func New{{.name}}() *{{.name}} {	
	return &{{.name}}{}
}

func (m *{{.name}})Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need 
		next(w, r)
	}	
}
`

func genMiddleware(dir string, middlewares []string) error {
	for _, item := range middlewares {
		filename := strings.TrimSuffix(strings.ToLower(item), "middleware") + "middleware" + ".go"
		fp, created, err := util.MaybeCreateFile(dir, middlewareDir, filename)
		if err != nil {
			return err
		}
		if !created {
			return nil
		}
		defer fp.Close()

		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		t := template.Must(template.New("contextTemplate").Parse(middlewareImplementCode))
		buffer := new(bytes.Buffer)
		err = t.Execute(buffer, map[string]string{
			"name": strings.Title(name),
		})
		if err != nil {
			return err
		}

		formatCode := formatCode(buffer.String())
		_, err = fp.WriteString(formatCode)
		return err
	}
	return nil
}
