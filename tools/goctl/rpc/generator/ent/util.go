package ent

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
)

func convertTypeToProtoType(typeName string) string {
	switch typeName {
	case "float32":
		typeName = "float"
	case "float64":
		typeName = "float"
	case "int":
		typeName = "int64"
	}
	return typeName
}

// convertSpecificNounToUpper is used to convert snack format to Ent format
func convertSpecificNounToUpper(str string) string {
	target := parser.CamelCase(str)
	target = strings.Replace(target, "Uuid", "UUID", -1)
	target = strings.Replace(target, "Api", "API", -1)
	target = strings.Replace(target, "Id", "ID", -1)

	return target
}
