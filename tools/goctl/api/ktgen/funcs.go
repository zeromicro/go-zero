package ktgen

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
)

var funcsMap = template.FuncMap{
	"lowCamelCase":    lowCamelCase,
	"routeToFuncName": routeToFuncName,
	"parseType":       parseType,
	"add":             add,
	"upperCase":       upperCase,
	"parseOptional":   parseOptional,
}

func lowCamelCase(s string) string {
	if len(s) < 1 {
		return ""
	}
	s = util.ToCamelCase(util.ToSnakeCase(s))
	return util.ToLower(s[:1]) + s[1:]
}

func routeToFuncName(method, path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.ReplaceAll(path, ":", "With_")

	return strings.ToLower(method) + strcase.ToCamel(path)
}

func parseType(t string) string {
	t = strings.Replace(t, "*", "", -1)
	if strings.HasPrefix(t, "[]") {
		return "List<" + parseType(t[2:]) + ">"
	}

	if strings.HasPrefix(t, "map") {
		tys, e := decomposeType(t)
		if e != nil {
			log.Fatal(e)
		}
		if len(tys) != 2 {
			log.Fatal("Map type number !=2")
		}
		return "Map<String," + parseType(tys[1]) + ">"
	}

	switch t {
	case "string":
		return "String"
	case "int", "int32", "int64":
		return "Int"
	case "float", "float32", "float64":
		return "Double"
	case "bool":
		return "Boolean"
	default:
		return t
	}
}

func decomposeType(t string) (result []string, err error) {
	add := func(tp string) error {
		ret, err := decomposeType(tp)
		if err != nil {
			return err
		}

		result = append(result, ret...)
		return nil
	}
	if strings.HasPrefix(t, "map") {
		t = strings.ReplaceAll(t, "map", "")
		if t[0] == '[' {
			pos := strings.Index(t, "]")
			if pos > 1 {
				if err = add(t[1:pos]); err != nil {
					return
				}
				if len(t) > pos+1 {
					err = add(t[pos+1:])
					return
				}
			}
		}
	} else if strings.HasPrefix(t, "[]") {
		if len(t) > 2 {
			err = add(t[2:])
			return
		}
	} else if strings.HasPrefix(t, "*") {
		err = add(t[1:])
		return
	} else {
		result = append(result, t)
		return
	}

	err = fmt.Errorf("bad type %q", t)
	return
}

func add(a, i int) int {
	return a + i
}

func upperCase(s string) string {
	return strings.ToUpper(s)
}

func parseOptional(m spec.Member) string {
	if m.IsOptionalOrOmitEmpty() {
		return "?"
	}
	return ""
}
