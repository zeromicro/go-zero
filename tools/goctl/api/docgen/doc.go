package docgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

//go:embed markdown.tpl
var markdownTemplate string

func genDoc(api *spec.ApiSpec, dir, filename string) error {
	if len(api.Service.Routes()) == 0 {
		return nil
	}

	fp, _, err := apiutil.MaybeCreateFile(dir, "", filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	var builder strings.Builder
	for index, route := range api.Service.Routes() {
		routeComment := route.JoinedDoc()
		if len(routeComment) == 0 {
			routeComment = "N/A"
		}

		requestContent, err := buildDoc(route.RequestType, api.Types)
		if err != nil {
			return err
		}

		responseContent, err := buildDoc(route.ResponseType, api.Types)
		if err != nil {
			return err
		}

		t := template.Must(template.New("markdownTemplate").Parse(markdownTemplate))
		var tmplBytes bytes.Buffer
		err = t.Execute(&tmplBytes, map[string]string{
			"index":           strconv.Itoa(index + 1),
			"routeComment":    routeComment,
			"method":          strings.ToUpper(route.Method),
			"uri":             route.Path,
			"requestType":     "`" + stringx.TakeOne(route.RequestTypeName(), "-") + "`",
			"responseType":    "`" + stringx.TakeOne(route.ResponseTypeName(), "-") + "`",
			"requestContent":  requestContent,
			"responseContent": responseContent,
		})
		if err != nil {
			return err
		}

		builder.Write(tmplBytes.Bytes())
	}

	_, err = fp.WriteString(strings.Replace(builder.String(), "&#34;", `"`, -1))
	return err
}

func buildDoc(route spec.Type, types []spec.Type) (string, error) {
	if route == nil || len(route.Name()) == 0 {
		return "", nil
	}

	tps := make([]spec.Type, 0)
	tps = append(tps, route)
	if definedType, ok := route.(spec.DefineStruct); ok {
		associatedTypes(definedType, &tps)
	}
	value, err := buildTypes(tps, types)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\n\n```golang\n%s\n```\n", value), nil
}

func associatedTypes(tp spec.DefineStruct, tps *[]spec.Type) {
	hasAdded := false
	for _, item := range *tps {
		if item.Name() == tp.Name() {
			hasAdded = true
			break
		}
	}
	if !hasAdded {
		*tps = append(*tps, tp)
	}

	for _, item := range tp.Members {
		if definedType, ok := item.Type.(spec.DefineStruct); ok {
			associatedTypes(definedType, tps)
		}
	}
}

// buildTypes gen types to string
func buildTypes(types, all []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp, all); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name()+" generate error")
		}
	}

	return builder.String(), nil
}

func writeType(writer io.Writer, tp spec.Type, all []spec.Type) error {
	fmt.Fprintf(writer, "type %s struct {\n", util.Title(tp.Name()))
	if err := writerMembers(writer, tp, all); err != nil {
		return err
	}
	fmt.Fprintf(writer, "}")
	return nil
}

func writerMembers(writer io.Writer, tp spec.Type, all []spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	getTargetType := func(tp string) spec.Type {
		for _, v := range all {
			if v.Name() == tp {
				return v
			}
		}
		return nil
	}
	for _, member := range structType.Members {
		if member.IsInline {
			inlineType := getTargetType(member.Type.Name())
			if inlineType == nil {
				if _, err := fmt.Fprintf(writer, "%s\n", strings.Title(member.Type.Name())); err != nil {
					return err
				}
			} else {
				if err := writerMembers(writer, inlineType, all); err != nil {
					return err
				}
			}

			continue
		}

		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type, 1); err != nil {
			return err
		}
	}

	return nil
}

func writeProperty(writer io.Writer, name, tag, comment string, tp spec.Type, indent int) error {
	apiutil.WriteIndent(writer, indent)
	var err error
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
		_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp.Name(), tag, comment)
	} else {
		_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp.Name(), tag)
	}

	return err
}
