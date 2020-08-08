package javagen

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const getSetTemplate = `
{{.indent}}{{.decorator}}
{{.indent}}public {{.returnType}} get{{.property}}() {
{{.indent}}	return this.{{.propertyValue}};
{{.indent}}}

{{.indent}}public void set{{.property}}({{.type}} {{.propertyValue}}) {
{{.indent}}	this.{{.propertyValue}} = {{.propertyValue}};
{{.indent}}}
`

func writeProperty(writer io.Writer, member spec.Member, indent int) error {
	writeIndent(writer, indent)
	ty, err := goTypeToJava(member.Type)
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
	writeDefaultValue(writer, member)
	fmt.Fprint(writer, ";\n")
	return err
}

func writeDefaultValue(writer io.Writer, member spec.Member) error {
	switch member.Type {
	case "string":
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

func writeBreakline(writer io.Writer) {
	fmt.Fprint(writer, "\n")
}

func isPrimitiveType(tp string) bool {
	switch tp {
	case "int", "int32", "int64":
		return true
	case "float", "float32", "float64":
		return true
	case "bool":
		return true
	}
	return false
}

func goTypeToJava(tp string) (string, error) {
	if len(tp) == 0 {
		return "", errors.New("property type empty")
	}
	if strings.HasPrefix(tp, "*") {
		tp = tp[1:]
	}
	switch tp {
	case "string":
		return "String", nil
	case "int64":
		return "long", nil
	case "int", "int8", "int32":
		return "int", nil
	case "float", "float32", "float64":
		return "double", nil
	case "bool":
		return "boolean", nil
	}
	if strings.HasPrefix(tp, "[]") {
		tys, err := apiutil.DecomposeType(tp)
		if err != nil {
			return "", err
		}
		if len(tys) == 0 {
			return "", fmt.Errorf("%s tp parse error", tp)
		}
		return fmt.Sprintf("java.util.ArrayList<%s>", util.Title(tys[0])), nil
	} else if strings.HasPrefix(tp, "map") {
		tys, err := apiutil.DecomposeType(tp)
		if err != nil {
			return "", err
		}
		if len(tys) == 2 {
			return "", fmt.Errorf("%s tp parse error", tp)
		}
		return fmt.Sprintf("java.util.HashMap<String, %s>", util.Title(tys[1])), nil
	}
	return util.Title(tp), nil
}

func genGetSet(writer io.Writer, tp spec.Type, indent int) error {
	t := template.Must(template.New("getSetTemplate").Parse(getSetTemplate))
	for _, member := range tp.Members {
		var tmplBytes bytes.Buffer

		oty, err := goTypeToJava(member.Type)
		if err != nil {
			return err
		}
		tyString := oty
		decorator := ""
		if !isPrimitiveType(member.Type) {
			if member.IsOptional() {
				decorator = "@org.jetbrains.annotations.Nullable "
			} else {
				decorator = "@org.jetbrains.annotations.NotNull "
			}
			tyString = decorator + tyString
		}

		err = t.Execute(&tmplBytes, map[string]string{
			"property":      util.Title(member.Name),
			"propertyValue": util.Untitle(member.Name),
			"type":          tyString,
			"decorator":     decorator,
			"returnType":    oty,
			"indent":        indentString(indent),
		})
		if err != nil {
			return err
		}
		r := tmplBytes.String()
		r = strings.Replace(r, " boolean get", " boolean is", 1)
		writer.Write([]byte(r))
	}
	return nil
}
