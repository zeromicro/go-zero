package gogen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/util"
	"github.com/lerity-yao/go-zero/tools/cztctl/pkg/golang"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
	"github.com/zeromicro/go-zero/core/collection"
)

type FileGenConfig struct {
	Dir             string
	Subdir          string
	Filename        string
	TemplateName    string
	Category        string
	TemplateFile    string
	BuiltinTemplate string
	Data            any
}

func GenFile(c FileGenConfig) error {
	fp, created, err := util.MaybeCreateFile(c.Dir, c.Subdir, c.Filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var text string
	if len(c.Category) == 0 || len(c.TemplateFile) == 0 {
		text = c.BuiltinTemplate
	} else {
		text, err = pathx.LoadTemplate(c.Category, c.TemplateFile, c.BuiltinTemplate)
		if err != nil {
			return err
		}
	}

	t := template.Must(template.New(c.TemplateName).Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, c.Data)
	if err != nil {
		return err
	}

	code := golang.FormatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}

func writeProperty(writer io.Writer, name, tag, comment string, tp spec.Type, indent int) error {
	util.WriteIndent(writer, indent)
	var (
		err            error
		isNestedStruct bool
	)
	structType, ok := tp.(spec.NestedStruct)
	if ok {
		isNestedStruct = true
	}
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
	}

	if isNestedStruct {
		_, err = fmt.Fprintf(writer, "%s struct {\n", strings.Title(name))
		if err != nil {
			return err
		}

		if err := writeMember(writer, structType.Members); err != nil {
			return err
		}

		_, err := fmt.Fprintf(writer, "} %s", tag)
		if err != nil {
			return err
		}

		if len(comment) > 0 {
			_, err = fmt.Fprintf(writer, " %s", comment)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprint(writer, "\n")
		if err != nil {
			return err
		}
	} else {
		if len(comment) > 0 {
			_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp.Name(), tag, comment)
			if err != nil {
				return err
			}
		} else {
			_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp.Name(), tag)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getAuths(api *spec.ApiSpec) []string {
	authNames := collection.NewSet[string]()
	for _, g := range api.Service.Groups {
		jwt := g.GetAnnotation("jwt")
		if len(jwt) > 0 {
			authNames.Add(jwt)
		}
	}
	return authNames.Keys()
}

func getJwtTrans(api *spec.ApiSpec) []string {
	jwtTransList := collection.NewSet[string]()
	for _, g := range api.Service.Groups {
		jt := g.GetAnnotation(jwtTransKey)
		if len(jt) > 0 {
			jwtTransList.Add(jt)
		}
	}
	return jwtTransList.Keys()
}

func getMiddleware(api *spec.ApiSpec) []string {
	result := collection.NewSet[string]()
	for _, g := range api.Service.Groups {
		middleware := g.GetAnnotation("middleware")
		if len(middleware) > 0 {
			for _, item := range strings.Split(middleware, ",") {
				result.Add(strings.TrimSpace(item))
			}
		}
	}

	return result.Keys()
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

func getDoc(doc string) string {
	if len(doc) == 0 {
		return ""
	}

	return "// " + strings.Trim(doc, "\"")
}
