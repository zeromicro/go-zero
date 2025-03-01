package gen

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen/template"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func IsOptionalOrOmitEmpty(m spec.Member) bool {
	tag := m.Tags()
	for _, item := range tag {
		if stringx.Contains(item.Options, "optional") || stringx.Contains(item.Options, "omitempty") {
			return true
		}
	}
	return false
}

func GenTsType(m spec.Member, indent int) (ty string, err error) {
	v, ok := m.Type.(spec.NestedStruct)
	if ok {
		nestedIndent := indent + 4
		ctt := &template.ComponentNestedTypeTemplateData{
			Indent:   nestedIndent,
			TypeName: m.Type.Name(),
		}
		writer := bytes.NewBuffer(nil)
		if ms, err := BuildMembers(v, false, nestedIndent); err != nil {
			return "", err
		} else {
			ctt.Members = ms
		}
		if err := template.GenTs(writer, template.Nested, ctt); err != nil {
			return "", err
		}
		return writer.String(), nil
	}

	ty, err = GoTypeToTs(m.Type, false)
	if enums := m.GetEnumOptions(); enums != nil {
		if ty == "string" {
			for i := range enums {
				enums[i] = "'" + enums[i] + "'"
			}
		}
		ty = strings.Join(enums, " | ")
	}
	return
}

func GoTypeToTs(tp spec.Type, fromPacket bool) (string, error) {
	switch v := tp.(type) {
	case spec.DefineStruct:
		return addPrefix(tp, fromPacket), nil
	case spec.PrimitiveType:
		r, ok := primitiveType(tp.Name())
		if !ok {
			return "", errors.New("unsupported primitive type " + tp.Name())
		}

		return r, nil
	case spec.MapType:
		valueType, err := GoTypeToTs(v.Value, fromPacket)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("{ [key: string]: %s }", valueType), nil
	case spec.ArrayType:
		if tp.Name() == "[]byte" {
			return "Blob", nil
		}

		valueType, err := GoTypeToTs(v.Value, fromPacket)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Array<%s>", valueType), nil
	case spec.InterfaceType:
		return "any", nil
	case spec.PointerType:
		return GoTypeToTs(v.Type, fromPacket)
	}

	return "", errors.New("unsupported type " + tp.Name())
}

func addPrefix(tp spec.Type, fromPacket bool) string {
	if fromPacket {
		return packagePrefix + util.Title(tp.Name())
	}
	return util.Title(tp.Name())
}

func primitiveType(tp string) (string, bool) {
	switch tp {
	case "string":
		return "string", true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number", true
	case "float", "float32", "float64":
		return "number", true
	case "bool":
		return "boolean", true
	case "[]byte":
		return "Blob", true
	case "interface{}":
		return "any", true
	}
	return "", false
}
