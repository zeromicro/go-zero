package dartgen

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
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

func getBaseName(str string) string {
	return path.Base(str)
}

func getPropertyFromMember(member spec.Member) string {
	name, err := member.GetPropertyName()
	if err != nil {
		panic(fmt.Sprintf("cannot get property name of %q", member.Name))
	}
	return name
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

func buildSpecType(tp spec.Type, name string) spec.Type {
	switch v := tp.(type) {
	case spec.PrimitiveType:
		return spec.PrimitiveType{RawName: name}
	case spec.MapType:
		return spec.MapType{RawName: name, Key: v.Key, Value: v.Value}
	case spec.ArrayType:
		return spec.ArrayType{RawName: name, Value: v.Value}
	case spec.InterfaceType:
		return spec.InterfaceType{RawName: name}
	case spec.PointerType:
		return spec.PointerType{RawName: name, Type: v.Type}
	}
	return tp
}

func specTypeToDart(tp spec.Type) (string, error) {
	switch v := tp.(type) {
	case spec.DefineStruct:
		return tp.Name(), nil
	case spec.PrimitiveType:
		r, ok := primitiveType(tp.Name())
		if !ok {
			return "", errors.New("unsupported primitive type " + tp.Name())
		}
		return r, nil
	case spec.MapType:
		valueType, err := specTypeToDart(v.Value)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Map<String, %s>", valueType), nil
	case spec.ArrayType:
		if tp.Name() == "[]byte" {
			return "List<int>", nil
		}

		valueType, err := specTypeToDart(v.Value)
		if err != nil {
			return "", err
		}

		s := getBaseType(valueType)
		if len(s) != 0 {
			return s, nil
		}
		return fmt.Sprintf("List<%s>", valueType), nil
	case spec.InterfaceType:
		return "Object", nil
	case spec.PointerType:
		return specTypeToDart(v.Type)
	}

	return "", errors.New("unsupported primitive type " + tp.Name())
}

func getBaseType(valueType string) string {
	switch valueType {
	case "int":
		return "List<int>"
	case "double":
		return "List<double>"
	case "boolean":
		return "List<bool>"
	case "String":
		return "List<String>"
	default:
		return ""
	}
}

func primitiveType(tp string) (string, bool) {
	switch tp {
	case "string":
		return "String", true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "rune":
		return "int", true
	case "float32", "float64":
		return "double", true
	case "bool":
		return "bool", true
	}

	return "", false
}
