package javagen

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func writeProperty(writer io.Writer, member spec.Member, indent int) error {
	if len(member.Comment) > 0 {
		writeIndent(writer, indent)
		fmt.Fprint(writer, member.Comment+util.NL)
	}
	writeIndent(writer, indent)
	ty, err := specTypeToJava(member.Type)
	ty = strings.Replace(ty, "*", "", 1)
	if err != nil {
		return err
	}

	name, err := member.GetPropertyName()
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(writer, "private %s %s", ty, name)
	if err != nil {
		return err
	}

	err = writeDefaultValue(writer, member)
	if err != nil {
		return err
	}

	fmt.Fprint(writer, ";\n")
	return err
}

func writeDefaultValue(writer io.Writer, member spec.Member) error {
	javaType, err := specTypeToJava(member.Type)
	if err != nil {
		return err
	}

	if javaType == "String" {
		_, err := fmt.Fprintf(writer, " = \"\"")
		return err
	}
	return nil
}

func writeIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

func indentString(indent int) string {
	var result = ""
	for i := 0; i < indent; i++ {
		result += "\t"
	}
	return result
}

func specTypeToJava(tp spec.Type) (string, error) {
	switch v := tp.(type) {
	case spec.DefineStruct:
		return util.Title(tp.Name()), nil
	case spec.PrimitiveType:
		r, ok := primitiveType(tp.Name())
		if !ok {
			return "", errors.New("unsupported primitive type " + tp.Name())
		}
		return r, nil
	case spec.MapType:
		valueType, err := specTypeToJava(v.Value)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("java.util.HashMap<String, %s>", util.Title(valueType)), nil
	case spec.ArrayType:
		if tp.Name() == "[]byte" {
			return "byte[]", nil
		}

		valueType, err := specTypeToJava(v.Value)
		if err != nil {
			return "", err
		}

		switch valueType {
		case "int":
			return "Integer[]", nil
		case "long":
			return "Long[]", nil
		case "float":
			return "Float[]", nil
		case "double":
			return "Double[]", nil
		case "boolean":
			return "Boolean[]", nil
		}

		return fmt.Sprintf("java.util.ArrayList<%s>", util.Title(valueType)), nil
	case spec.InterfaceType:
		return "Object", nil
	case spec.PointerType:
		return specTypeToJava(v.Type)
	}

	return "", errors.New("unsupported primitive type " + tp.Name())
}

func primitiveType(tp string) (string, bool) {
	switch tp {
	case "string":
		return "String", true
	case "int64", "uint64":
		return "long", true
	case "int", "int8", "int32", "uint", "uint8", "uint16", "uint32":
		return "int", true
	case "float", "float32":
		return "float", true
	case "float64":
		return "double", true
	case "bool":
		return "boolean", true
	}

	return "", false
}
