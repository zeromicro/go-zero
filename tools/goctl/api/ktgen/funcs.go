package ktgen

import (
	"log"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
)

var funcsMap = template.FuncMap{
	"lowCamelCase":    lowCamelCase,
	"routeToFuncName": routeToFuncName,
	"parseType":       parseType,
	"add":             add,
	"upperCase":       upperCase,
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
		tys, e := util.DecomposeType(t)
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

func add(a, i int) int {
	return a + i
}

func upperCase(s string) string {
	return strings.ToUpper(s)
}
