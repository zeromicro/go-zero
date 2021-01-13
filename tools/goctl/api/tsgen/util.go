package tsgen

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func writeProperty(writer io.Writer, member spec.Member, indent int) error {
	writeIndent(writer, indent)
	ty, err := goTypeToTs(member.Type, false)
	if err != nil {
		return err
	}

	optionalTag := ""
	if member.IsOptional() || member.IsOmitempty() {
		optionalTag = "?"
	}
	name, err := member.GetPropertyName()
	if err != nil {
		return err
	}

	comment := member.GetComment()
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = " // " + strings.TrimSpace(comment)
	}
	if len(member.Docs) > 0 {
		fmt.Fprintf(writer, "%s\n", strings.Join(member.Docs, ""))
		writeIndent(writer, 1)
	}
	_, err = fmt.Fprintf(writer, "%s%s: %s%s\n", name, optionalTag, ty, comment)
	return err
}

func writeIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

func goTypeToTs(tp spec.Type, fromPacket bool) (string, error) {
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
		valueType, err := goTypeToTs(v.Value, fromPacket)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("{ [key: string]: %s }", valueType), nil
	case spec.ArrayType:
		if tp.Name() == "[]byte" {
			return "Blob", nil
		}

		valueType, err := goTypeToTs(v.Value, fromPacket)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Array<%s>", valueType), nil
	case spec.InterfaceType:
		return "any", nil
	case spec.PointerType:
		return goTypeToTs(v.Type, fromPacket)
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
	case "int", "int8", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
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

func writeType(writer io.Writer, tp spec.Type) error {
	fmt.Fprintf(writer, "export interface %s {\n", util.Title(tp.Name()))
	if err := writeMembers(writer, tp, false); err != nil {
		return err
	}

	fmt.Fprintf(writer, "}\n")
	return genParamsTypesIfNeed(writer, tp)
}

func genParamsTypesIfNeed(writer io.Writer, tp spec.Type) error {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		return errors.New("no members of type " + tp.Name())
	}

	members := definedType.GetNonBodyMembers()
	if len(members) == 0 {
		return nil
	}
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "export interface %sParams {\n", util.Title(tp.Name()))
	if err := writeMembers(writer, tp, true); err != nil {
		return err
	}

	fmt.Fprintf(writer, "}\n")
	return nil
}

func writeMembers(writer io.Writer, tp spec.Type, isParam bool) error {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return writeMembers(writer, pointType.Type, isParam)
		}
		return errors.New(fmt.Sprintf("type %s not supported", tp.Name()))
	}

	members := definedType.GetBodyMembers()
	if isParam {
		members = definedType.GetNonBodyMembers()
	}
	for _, member := range members {
		if member.IsInline {
			if err := writeMembers(writer, member.Type, isParam); err != nil {
				return err
			}
			continue
		}

		if err := writeProperty(writer, member, 1); err != nil {
			return apiutil.WrapErr(err, " type "+tp.Name())
		}
	}
	return nil
}
