package docgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
)

//go:embed markdown.tpl
var markdownTemplate string

func genDoc(api *spec.ApiSpec, dir, filename string) error {
	if len(api.Service.Routes()) == 0 {
		return nil
	}

	fp, _, err := util.MaybeCreateFile(dir, "", filename)
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

		requestContent, err := buildDoc(route.RequestType)
		if err != nil {
			return err
		}

		responseContent, err := buildDoc(route.ResponseType)
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

func buildDoc(route spec.Type) (string, error) {
	if route == nil || len(route.Name()) == 0 {
		return "", nil
	}

	tps := make([]spec.Type, 0)
	tps = append(tps, route)
	if definedType, ok := route.(spec.DefineStruct); ok {
		associatedTypes(definedType, &tps)
	}
	value, err := gogen.BuildTypes(tps)
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
