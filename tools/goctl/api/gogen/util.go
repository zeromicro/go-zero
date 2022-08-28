package gogen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type fileGenConfig struct {
	dir             string
	subdir          string
	filename        string
	templateName    string
	category        string
	templateFile    string
	builtinTemplate string
	data            interface{}
}

func genFile(c fileGenConfig) error {
	fp, created, err := util.MaybeCreateFile(c.dir, c.subdir, c.filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var text string
	if len(c.category) == 0 || len(c.templateFile) == 0 {
		text = c.builtinTemplate
	} else {
		text, err = pathx.LoadTemplate(c.category, c.templateFile, c.builtinTemplate)
		if err != nil {
			return err
		}
	}

	t := template.Must(template.New(c.templateName).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, c.data)
	if err != nil {
		return err
	}

	code := golang.FormatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}

func writeProperty(writer io.Writer, name, tag, comment string, tp spec.Type, indent int, api *spec.ApiSpec) error {
	util.WriteIndent(writer, indent)
	var err error
	var refPropertyName = tp.Name()
	if isCustomType(refPropertyName) {
		strs := getRefProperty(api, refPropertyName, name)
		_, err = fmt.Fprintf(writer, "%s\n", strs)
	} else {
		if len(comment) > 0 {
			comment = strings.TrimPrefix(comment, "//")
			comment = "//" + comment
			_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp.Name(), tag, comment)
		} else {
			_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp.Name(), tag)
		}
	}

	return err
}

func getAuths(api *spec.ApiSpec) []string {
	authNames := collection.NewSet()
	for _, g := range api.Service.Groups {
		jwt := g.GetAnnotation("jwt")
		if len(jwt) > 0 {
			authNames.Add(jwt)
		}
	}
	return authNames.KeysStr()
}

func getJwtTrans(api *spec.ApiSpec) []string {
	jwtTransList := collection.NewSet()
	for _, g := range api.Service.Groups {
		jt := g.GetAnnotation(jwtTransKey)
		if len(jt) > 0 {
			jwtTransList.Add(jt)
		}
	}
	return jwtTransList.KeysStr()
}

func getMiddleware(api *spec.ApiSpec) []string {
	result := collection.NewSet()
	for _, g := range api.Service.Groups {
		middleware := g.GetAnnotation("middleware")
		if len(middleware) > 0 {
			for _, item := range strings.Split(middleware, ",") {
				result.Add(strings.TrimSpace(item))
			}
		}
	}

	return result.KeysStr()
}

func responseGoTypeName(r spec.Route, pkg ...string) string {
	if r.ResponseType == nil {
		return ""
	}

	resp := golangExpr(r.ResponseType, pkg...)
	switch r.ResponseType.(type) {
	case spec.DefineStruct:
		if !strings.HasPrefix(resp, "*") {
			return "*" + resp
		}
	}

	return resp
}

func requestGoTypeName(r spec.Route, pkg ...string) string {
	if r.RequestType == nil {
		return ""
	}

	return golangExpr(r.RequestType, pkg...)
}

func golangExpr(ty spec.Type, pkg ...string) string {
	switch v := ty.(type) {
	case spec.PrimitiveType:
		return v.RawName
	case spec.DefineStruct:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("%s.%s", pkg[0], strings.Title(v.RawName))
	case spec.ArrayType:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("[]%s", golangExpr(v.Value, pkg...))
	case spec.MapType:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("map[%s]%s", v.Key, golangExpr(v.Value, pkg...))
	case spec.PointerType:
		if len(pkg) > 1 {
			panic("package cannot be more than 1")
		}

		if len(pkg) == 0 {
			return v.RawName
		}

		return fmt.Sprintf("*%s", golangExpr(v.Type, pkg...))
	case spec.InterfaceType:
		return v.RawName
	}

	return ""
}

func isCustomType(t string) bool {
	var builtinType = []string{"string", "bool", "int", "uint", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "uintptr", "complex64", "complex128"}
	var is bool = true
	for _, v := range builtinType {
		if t == v {
			is = false
			break
		}
	}
	return is
}

// Generate nested types recursively
func getRefProperty(api *spec.ApiSpec, refPropertyName string, name string) string {
	var str string = ""
	for _, t := range api.Types {
		if strings.TrimLeft(refPropertyName, "*") == t.Name() {
			switch tm := t.(type) {
			case spec.DefineStruct:
				for _, m := range tm.Members {
					if isCustomType(m.Type.Name()) {
						// recursive
						str += getRefProperty(api, m.Type.Name(), m.Name)
					} else {
						if len(m.Comment) > 0 {
							comment := strings.TrimPrefix(m.Comment, "//")
							comment = "//" + comment
							str += fmt.Sprintf("%s %s %s %s\n\t", m.Name, m.Type.Name(), m.Tag, comment)
						} else {
							str += fmt.Sprintf("%s %s %s\n\t", m.Name, m.Type.Name(), m.Tag)
						}

					}

				}
			}
		}
	}
	if name == "" {
		temp := `${str}`
		return os.Expand(temp, func(k string) string {
			return str
		})
	} else {
		temp := `${name} struct {
			${str}}`
		return os.Expand(temp, func(k string) string {
			return map[string]string{
				"name": name,
				"str":  str,
			}[k]
		})
	}
}
