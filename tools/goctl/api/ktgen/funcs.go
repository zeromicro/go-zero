package ktgen

import (
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"log"
	"strings"
	"text/template"
)
var funcsMap=template.FuncMap{
	"lowCamelCase":lowCamelCase,
	"pathToFuncName":pathToFuncName,
	"parseType":parseType,
	"add":add,
}
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

	path = strings.Replace(path, "/", "_", -1)
	path = strings.Replace(path, "-", "_", -1)

	camel := util.ToCamelCase(path)
	return util.ToLower(camel[:1]) + camel[1:]
}
func parseType(t string) string {
	t=strings.Replace(t,"*","",-1)
	if strings.HasPrefix(t,"[]"){
		return "List<"+parseType(t[2:])+ ">"
	}

	if strings.HasPrefix(t,"map"){
		tys,e:=util.DecomposeType(t)
		if e!=nil{
		    log.Fatal(e)
		}
		if len(tys)!=2{
			log.Fatal("Map type number !=2")
		}
		return "Map<String,"+parseType(tys[1])+">"
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

func add(a,i int)int{
	return a+i
}