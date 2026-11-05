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

const (
	formTagKey   = "form"
	pathTagKey   = "path"
	headerTagKey = "header"
)

func normalizeHandlerName(handlerName string) string {
	handler := strings.Replace(handlerName, "Handler", "", 1)
	handler = lowCamelCase(handler)
	return handler
}

func lowCamelCase(s string) string {
	if len(s) < 1 {
		return ""
	}

	s = util.ToCamelCase(util.ToSnakeCase(s))
	return util.ToLower(s[:1]) + s[1:]
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

type promotedMember struct {
	member spec.Member
	depth  int
}

func getMembers(tp spec.Type) ([]spec.Member, error) {
	members, err := collectMembers(tp, make(map[string]bool), 0)
	if err != nil {
		return nil, err
	}

	selected := make(map[string]int)
	ambiguous := make(map[string]bool)
	for i, candidate := range members {
		name := lowCamelCase(candidate.member.Name)
		previous, ok := selected[name]
		if ok {
			previousProperty, err := members[previous].member.GetPropertyName()
			if err != nil {
				return nil, err
			}
			property, err := candidate.member.GetPropertyName()
			if err != nil {
				return nil, err
			}
			if previousProperty != property {
				return nil, fmt.Errorf("conflicting property names %q and %q for Dart field %s",
					previousProperty, property, name)
			}
		}
		if !ok || candidate.depth < members[previous].depth {
			selected[name] = i
			ambiguous[name] = false
		} else if candidate.depth == members[previous].depth {
			ambiguous[name] = true
		}
	}

	var result []spec.Member
	for i, candidate := range members {
		name := lowCamelCase(candidate.member.Name)
		if selected[name] != i {
			continue
		}
		if ambiguous[name] {
			return nil, fmt.Errorf("ambiguous embedded field %s", candidate.member.Name)
		}
		result = append(result, candidate.member)
	}
	return result, nil
}

func collectMembers(tp spec.Type, resolving map[string]bool, depth int) ([]promotedMember, error) {
	var members []spec.Member
	switch v := tp.(type) {
	case spec.DefineStruct:
		members = v.Members
	case spec.NestedStruct:
		members = v.Members
	case spec.PointerType:
		return collectMembers(v.Type, resolving, depth)
	default:
		return nil, fmt.Errorf("type %s is not a struct", tp.Name())
	}

	if resolving[tp.Name()] {
		return nil, fmt.Errorf("cyclic embedded type %s", tp.Name())
	}
	resolving[tp.Name()] = true
	defer delete(resolving, tp.Name())

	var result []promotedMember
	for _, member := range members {
		if member.IsInline {
			embedded, err := collectMembers(member.Type, resolving, depth+1)
			if err != nil {
				return nil, err
			}
			result = append(result, embedded...)
		} else {
			result = append(result, promotedMember{member: member, depth: depth})
		}
	}
	return result, nil
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

func isNumberType(s string) bool {
	switch s {
	case "int", "double":
		return true
	default:
		return false
	}
}

func isListType(s string) bool {
	return strings.HasPrefix(s, "List<")
}

func isClassListType(s string) bool {
	return isListType(s) && !isAtomicType(getCoreType(s))
}

func isAtomicListType(s string) bool {
	return isListType(s) && isAtomicType(getCoreType(s))
}

func isListItemsNullable(s string) bool {
	return isListType(s) && isNullableType(getCoreType(s))
}

func isMapType(s string) bool {
	return strings.HasPrefix(s, "Map<")
}

// Only interface types are nullable
func isNullableType(s string) bool {
	return strings.HasSuffix(s, "?")
}

func appendNullCoalescing(member spec.Member) string {
	if isNullableType(member.Type.Name()) {
		return "m['" + getPropertyFromMember(member) + "'] == null ? null : "
	}
	return ""
}

// To be compatible with omitempty tags in Golang
// Only set default value for non-nullable types
func appendDefaultEmptyValue(s string) string {
	if isNullableType(s) {
		return ""
	}

	if isAtomicType(s) {
		switch s {
		case "String":
			return `?? ""`
		case "int":
			return "?? 0"
		case "double":
			return "?? 0.0"
		case "bool":
			return "?? false"
		default:
			panic(errors.New("unknown atomic type"))
		}
	}
	if isListType(s) {
		return "?? []"
	}
	if isMapType(s) {
		return "?? {}"
	}
	return ""
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
		return "Object?", nil
	case spec.PointerType:
		valueType, err := specTypeToDart(v.Type)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s?", valueType), nil
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

func hasUrlPathParams(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(route.RequestTypeName()) > 0 && len(ds.GetTagMembers(pathTagKey)) > 0
}

func extractPositionalParamsFromPath(route spec.Route) (string, error) {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return "", nil
	}

	members, err := getMembers(ds)
	if err != nil {
		return "", err
	}
	params := make([]string, 0, len(members))
	for _, member := range members {
		if !member.IsTagMember(pathTagKey) {
			continue
		}
		dartType, err := specTypeToDart(member.Type)
		if err != nil {
			return "", err
		}
		params = append(params, fmt.Sprintf("%s %s", dartType, getPropertyFromMember(member)))
	}

	return strings.Join(params, ", "), nil
}

func makeDartRequestUrlPath(route spec.Route) (string, error) {
	path := route.Path
	if route.RequestType == nil {
		return `"` + path + `"`, nil
	}

	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return path, nil
	}

	members, err := getMembers(ds)
	if err != nil {
		return "", err
	}
	for _, member := range members {
		if !member.IsTagMember(pathTagKey) {
			continue
		}
		paramName := member.Tags()[0].Name
		path = strings.ReplaceAll(path, ":"+paramName, "${"+getPropertyFromMember(member)+"}")
	}

	return `"` + path + `"`, nil
}
