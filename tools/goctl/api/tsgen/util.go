package tsgen

import (
	"fmt"
	"io"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func writeProperty(writer io.Writer, member spec.Member, indent int, prefixForType func(string) string) error {
	writeIndent(writer, indent)
	ty, err := goTypeToTs(member.Type, prefixForType)
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

func goTypeToTs(tp string, prefixForType func(string) string) (string, error) {
	if val, pri := primitiveType(tp); pri {
		return val, nil
	}
	if tp == "[]byte" {
		return "Blob", nil
	} else if strings.HasPrefix(tp, "[][]") {
		tys, err := apiutil.DecomposeType(tp)
		if err != nil {
			return "", err
		}
		if len(tys) == 0 {
			return "", fmt.Errorf("%s tp parse error", tp)
		}
		innerType, err := goTypeToTs(tys[0], prefixForType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Array<Array<%s>>", innerType), nil
	} else if strings.HasPrefix(tp, "[]") {
		tys, err := apiutil.DecomposeType(tp)
		if err != nil {
			return "", err
		}
		if len(tys) == 0 {
			return "", fmt.Errorf("%s tp parse error", tp)
		}
		innerType, err := goTypeToTs(tys[0], prefixForType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Array<%s>", innerType), nil
	} else if strings.HasPrefix(tp, "map") {
		tys, err := apiutil.DecomposeType(tp)
		if err != nil {
			return "", err
		}
		if len(tys) != 2 {
			return "", fmt.Errorf("%s tp parse error", tp)
		}
		innerType, err := goTypeToTs(tys[1], prefixForType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("{ [key: string]: %s }", innerType), nil
	}
	return addPrefixIfNeed(util.Title(tp), prefixForType), nil
}

func addPrefixIfNeed(tp string, prefixForType func(string) string) string {
	if val, pri := primitiveType(tp); pri {
		return val
	}
	tp = strings.Replace(tp, "*", "", 1)
	return prefixForType(tp) + util.Title(tp)
}

func primitiveType(tp string) (string, bool) {
	switch tp {
	case "string":
		return "string", true
	case "int", "int8", "int32", "int64":
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

func writeType(writer io.Writer, tp spec.Type, inlineType func(string) (*spec.Type, error), prefixForType func(string) string) error {
	fmt.Fprintf(writer, "export interface %s {\n", util.Title(tp.Name))
	if err := genMembers(writer, tp, false, inlineType, prefixForType); err != nil {
		return err
	}
	fmt.Fprintf(writer, "}\n")
	err := genParamsTypesIfNeed(writer, tp, inlineType, prefixForType)
	if err != nil {
		return err
	}
	return nil
}

func genParamsTypesIfNeed(writer io.Writer, tp spec.Type, inlineType func(string) (*spec.Type, error), prefixForType func(string) string) error {
	members := tp.GetNonBodyMembers()
	if len(members) == 0 {
		return nil
	}
	fmt.Fprintf(writer, "\n")
	fmt.Fprintf(writer, "export interface %sParams {\n", util.Title(tp.Name))
	if err := genMembers(writer, tp, true, inlineType, prefixForType); err != nil {
		return err
	}
	fmt.Fprintf(writer, "}\n")
	return nil
}

func genMembers(writer io.Writer, tp spec.Type, isParam bool, inlineType func(string) (*spec.Type, error), prefixForType func(string) string) error {
	members := tp.GetBodyMembers()
	if isParam {
		members = tp.GetNonBodyMembers()
	}
	for _, member := range members {
		if member.IsInline {
			// 获取inline类型的成员然后添加到type中
			it, err := inlineType(strings.TrimPrefix(member.Type, "*"))
			if err != nil {
				return err
			}
			if err := genMembers(writer, *it, isParam, inlineType, prefixForType); err != nil {
				return err
			}
			continue
		}
		if err := writeProperty(writer, member, 1, prefixForType); err != nil {
			return apiutil.WrapErr(err, " type "+tp.Name)
		}
	}
	return nil
}
