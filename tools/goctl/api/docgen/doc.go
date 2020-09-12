package docgen

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
)

const (
	markdownTemplate = `
### {{.index}}. {{.routeComment}}

1. 路由定义

- Url: {{.uri}}
- Method: {{.method}}
- Request: {{.requestType}}
- Response: {{.responseType}}


2. 类型定义 

{{.responseContent}}  

`
)

func genDoc(api *spec.ApiSpec, dir string, filename string) error {
	fp, _, err := util.MaybeCreateFile(dir, "", filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	var builder strings.Builder
	for index, route := range api.Service.Routes {
		routeComment, _ := util.GetAnnotationValue(route.Annotations, "doc", "summary")
		if len(routeComment) == 0 {
			routeComment = "N/A"
		}

		responseContent, err := responseBody(api, route)
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
			"requestType":     "`" + stringx.TakeOne(route.RequestType.Name, "-") + "`",
			"responseType":    "`" + stringx.TakeOne(route.ResponseType.Name, "-") + "`",
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

func responseBody(api *spec.ApiSpec, route spec.Route) (string, error) {
	tps := util.GetLocalTypes(api, route)
	value, err := gogen.BuildTypes(tps)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\n\n```golang\n%s\n```\n", value), nil
}
