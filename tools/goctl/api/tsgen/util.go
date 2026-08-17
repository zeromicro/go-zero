package tsgen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

const (
	formTagKey   = "form"
	pathTagKey   = "path"
	headerTagKey = "header"
)

func writeProperty(writer io.Writer, member spec.Member, indent int) error {
	writeIndent(writer, indent)
	ty, err := genTsType(member, indent)
	if err != nil {
		return err
	}

	optionalTag := ""
	if member.IsOptional() || member.IsOmitEmpty() {
		optionalTag = "?"
	}
	name, err := member.GetPropertyName()
	if err != nil {
		return err
	}

	if strings.Contains(name, "-") {
		name = fmt.Sprintf("'%s'", name)
	}

	comment := member.GetComment()
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = " // " + strings.TrimSpace(comment)
	}
	if len(member.Docs) > 0 {
		fmt.Fprintf(writer, "%s\n", strings.Join(member.Docs, ""))
		writeIndent(writer, indent)
	}
	_, err = fmt.Fprintf(writer, "%s%s: %s%s\n", name, optionalTag, ty, comment)
	return err
}

func writeIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

func genTsType(m spec.Member, indent int) (ty string, err error) {
	v, ok := m.Type.(spec.NestedStruct)
	if ok {
		writer := bytes.NewBuffer(nil)
		_, err := fmt.Fprintf(writer, "{\n")
		if err != nil {
			return "", err
		}

		if err := writeMembers(writer, v, false, indent+1); err != nil {
			return "", err
		}

		writeIndent(writer, indent)
		_, err = fmt.Fprintf(writer, "}")
		if err != nil {
			return "", err
		}
		return writer.String(), nil
	}

	ty, err = goTypeToTs(m.Type, false)
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

func writeType(writer io.Writer, tp spec.Type) error {
	fmt.Fprintf(writer, "export interface %s {\n", util.Title(tp.Name()))
	if err := writeMembers(writer, tp, false, 1); err != nil {
		return err
	}

	fmt.Fprintf(writer, "}\n")
	return genParamsTypesIfNeed(writer, tp)
}

func genParamsTypesIfNeed(writer io.Writer, tp spec.Type) error {
	_, ok := tp.(spec.DefineStruct)
	if !ok {
		return errors.New("no members of type " + tp.Name())
	}

	// Check if there are actual non-body members (recursively through inline structs)
	if !hasActualNonBodyMembers(tp) {
		return nil
	}

	fmt.Fprintf(writer, "export interface %sParams {\n", util.Title(tp.Name()))
	if err := writeTagMembers(writer, tp, formTagKey); err != nil {
		return err
	}
	fmt.Fprintf(writer, "}\n")

	if hasActualTagMembers(tp, headerTagKey) {
		fmt.Fprintf(writer, "export interface %sHeaders {\n", util.Title(tp.Name()))
		if err := writeTagMembers(writer, tp, headerTagKey); err != nil {
			return err
		}
		fmt.Fprintf(writer, "}\n")
	}

	return nil
}

func writeMembers(writer io.Writer, tp spec.Type, isParam bool, indent int) error {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return writeMembers(writer, pointType.Type, isParam, indent)
		}

		return fmt.Errorf("type %s not supported", tp.Name())
	}

	members := definedType.GetBodyMembers()
	if isParam {
		members = definedType.GetNonBodyMembers()
	}
	for _, member := range members {
		if member.IsInline {
			if err := writeMembers(writer, member.Type, isParam, indent); err != nil {
				return err
			}
			continue
		}

		if err := writeProperty(writer, member, indent); err != nil {
			return apiutil.WrapErr(err, " type "+tp.Name())
		}
	}
	return nil
}

func writeTagMembers(writer io.Writer, tp spec.Type, tagKey string) error {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return writeTagMembers(writer, pointType.Type, tagKey)
		}

		return fmt.Errorf("type %s not supported", tp.Name())
	}

	members := definedType.GetTagMembers(tagKey)
	for _, member := range members {
		if member.IsInline {
			if err := writeTagMembers(writer, member.Type, tagKey); err != nil {
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

// hasActualTagMembers checks if a type has actual members with the given tag,
// recursively checking inline/embedded structs
func hasActualTagMembers(tp spec.Type, tagKey string) bool {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return hasActualTagMembers(pointType.Type, tagKey)
		}
		return false
	}

	for _, m := range definedType.Members {
		if m.IsInline {
			// Recursively check inline members
			if hasActualTagMembers(m.Type, tagKey) {
				return true
			}
		} else {
			// Check non-inline members for the tag
			if m.IsTagMember(tagKey) {
				return true
			}
		}
	}
	return false
}

// hasActualBodyMembers checks if a type has actual body members (json tags),
// recursively checking inline/embedded structs
func hasActualBodyMembers(tp spec.Type) bool {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return hasActualBodyMembers(pointType.Type)
		}
		return false
	}

	for _, m := range definedType.Members {
		if m.IsInline {
			// Recursively check inline members
			if hasActualBodyMembers(m.Type) {
				return true
			}
		} else {
			// Check non-inline members for json tag
			if m.IsBodyMember() {
				return true
			}
		}
	}
	return false
}

// hasActualNonBodyMembers checks if a type has actual non-body members (form, path, header tags),
// recursively checking inline/embedded structs
func hasActualNonBodyMembers(tp spec.Type) bool {
	definedType, ok := tp.(spec.DefineStruct)
	if !ok {
		pointType, ok := tp.(spec.PointerType)
		if ok {
			return hasActualNonBodyMembers(pointType.Type)
		}
		return false
	}

	for _, m := range definedType.Members {
		if m.IsInline {
			// Recursively check inline members
			if hasActualNonBodyMembers(m.Type) {
				return true
			}
		} else {
			// Check non-inline members for non-body tags
			if !m.IsBodyMember() {
				return true
			}
		}
	}
	return false
}
