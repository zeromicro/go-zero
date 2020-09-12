package dartgen

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
)

func lowCamelCase(s string) string {
	if len(s) < 1 {
		return ""
	}
	s = util.ToCamelCase(util.ToSnakeCase(s))
	return util.ToLower(s[:1]) + s[1:]
}

func pathToFuncName(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(path, "/api") {
		path = "/api" + path
	}

	path = strings.Replace(path, "/", "_", -1)
	path = strings.Replace(path, "-", "_", -1)

	camel := util.ToCamelCase(path)
	return util.ToLower(camel[:1]) + camel[1:]
}

func tagGet(tag, k string) (reflect.Value, error) {
	v, _ := util.TagLookup(tag, k)
	out := strings.Split(v, ",")[0]
	return reflect.ValueOf(out), nil
}

func convertMemberType(api *spec.ApiSpec) {
	for i, t := range api.Types {
		for j, mem := range t.Members {
			api.Types[i].Members[j].Type = goTypeToDart(mem.Type)
		}
	}
}

func goTypeToDart(t string) string {
	t = strings.Replace(t, "*", "", -1)
	if strings.HasPrefix(t, "[]") {
		return "List<" + goTypeToDart(t[2:]) + ">"
	}

	if strings.HasPrefix(t, "map") {
		tys, e := util.DecomposeType(t)
		if e != nil {
			log.Fatal(e)
		}

		if len(tys) != 2 {
			log.Fatal("Map type number !=2")
		}

		return "Map<String," + goTypeToDart(tys[1]) + ">"
	}

	switch t {
	case "string":
		return "String"
	case "int", "int32", "int64":
		return "int"
	case "float", "float32", "float64":
		return "double"
	case "bool":
		return "bool"
	default:
		return t
	}
}

func isDirectType(s string) bool {
	return isAtomicType(s) || isListType(s) && isAtomicType(getCoreType(s))
}

func isAtomicType(s string) bool {
	switch s {
	case "String", "int", "double", "bool":
		return true
	default:
		return false
	}
}

func isListType(s string) bool {
	return strings.HasPrefix(s, "List<")
}

func isClassListType(s string) bool {
	return strings.HasPrefix(s, "List<") && !isAtomicType(getCoreType(s))
}

func getCoreType(s string) string {
	if isAtomicType(s) {
		return s
	}
	if isListType(s) {
		s = strings.Replace(s, "List<", "", -1)
		return strings.Replace(s, ">", "", -1)
	}
	return s
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
