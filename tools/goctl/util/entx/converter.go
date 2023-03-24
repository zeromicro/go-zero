// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	case "uint8", "uint16":
		typeName = "uint32"
	case "int8", "int16":
		typeName = "int32"
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
	case "uint8", "uint16":
		return "uint32"
	case "int8", "int16":
		return "int32"
	}
	return prop
}

// ConvertEntTypeToGotypeInSingleApi returns go type from ent type in single API service
func ConvertEntTypeToGotypeInSingleApi(prop string) string {
	switch prop {
	case "[16]byte":
		return "string"
	case "time.Time":
		return "int64"
	default:
		return prop
	}
}

// ConvertIDType returns uuid type by uuid flag
func ConvertIDType(useUUID bool) string {
	if useUUID {
		return "string"
	}
	return "uint64"
}

// ConvertOnlyEntTypeToGoType converts the type that only ent has to go type.
func ConvertOnlyEntTypeToGoType(t string) string {
	switch t {
	case "int8", "int16":
		return "int32"
	case "uint8", "uint16":
		return "uint32"
	default:
		return "uint32"
	}
}
