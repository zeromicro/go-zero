package entx

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

// ConvertEntTypeToProtoType returns prototype from ent type
func ConvertEntTypeToProtoType(typeName string) string {
	switch typeName {
	case "float32":
		typeName = "float"
	case "float64":
		typeName = "double"
	case "float":
		typeName = "double"
	case "int":
		typeName = "int64"
	case "uint":
		typeName = "uint64"
	case "[16]byte":
		typeName = "string"
	}
	return typeName
}

// ConvertProtoTypeToGoType returns go type from proto type
func ConvertProtoTypeToGoType(typeName string) string {
	switch typeName {
	case "float":
		typeName = "float32"
	case "double":
		typeName = "float64"
	}
	return typeName
}

// ConvertSpecificNounToUpper is used to convert snack format to Ent format
func ConvertSpecificNounToUpper(str string) string {
	target := parser.CamelCase(str)
	target = strings.Replace(target, "Uuid", "UUID", -1)
	target = strings.Replace(target, "Api", "API", -1)
	target = strings.Replace(target, "Id", "ID", -1)

	return target
}

// ConvertEntTypeToGotype returns go type from ent type
func ConvertEntTypeToGotype(prop string) string {
	switch prop {
	case "int":
		return "int64"
	case "uint":
		return "uint64"
	}
	return prop
}

// ConvertIDType returns uuid type by uuid flag
func ConvertIDType(useUUID bool) string {
	if useUUID {
		return "string"
	}
	return "uint64"
}
