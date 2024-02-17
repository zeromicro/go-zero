package dartgen

import (
	"os"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const apiTemplate = `import 'api.dart';
import '../data/{{with .Info}}{{getBaseName .Title}}{{end}}.dart';
{{with .Service}}
/// {{.Name}}
{{range .Routes}}
/// --{{.Path}}--
///
/// request: {{with .RequestType}}{{.Name}}{{end}}
/// response: {{with .ResponseType}}{{.Name}}{{end}}
Future {{pathToFuncName .Path}}( {{if ne .Method "get"}}{{with .RequestType}}{{.Name}} request,{{end}}{{end}}
    {Function({{with .ResponseType}}{{.Name}}{{end}}) ok,
    Function(String) fail,
    Function eventually}) async {
  await api{{if eq .Method "get"}}Get{{else}}Post{{end}}('{{.Path}}',{{if ne .Method "get"}}request,{{end}}
  	 ok: (data) {
    if (ok != null) ok({{with .ResponseType}}{{.Name}}.fromJson(data){{end}});
  }, fail: fail, eventually: eventually);
}
{{end}}
{{end}}`

const apiTemplateV2 = `import 'api.dart';
import '../data/{{with .Service}}{{.Name}}{{end}}.dart';
{{with .Service}}
/// {{.Name}}
{{range $i, $Route := .Routes}}
/// --{{.Path}}--
///
/// request: {{with .RequestType}}{{.Name}}{{end}}
/// response: {{with .ResponseType}}{{.Name}}{{end}}
Future {{normalizeHandlerName .Handler}}(
	{{if hasUrlPathParams $Route}}{{extractPositionalParamsFromPath $Route}},{{end}}
	{{if ne .Method "get"}}{{with .RequestType}}{{.Name}} request,{{end}}{{end}}
    {Function({{with .ResponseType}}{{.Name}}{{end}})? ok,
    Function(String)? fail,
    Function? eventually}) async {
  await api{{if eq .Method "get"}}Get{{else}}Post{{end}}({{makeDartRequestUrlPath $Route}},{{if ne .Method "get"}}request,{{end}}
  	 ok: (data) {
    if (ok != null) ok({{with .ResponseType}}{{.Name}}.fromJson(data){{end}});
  }, fail: fail, eventually: eventually);
}
{{end}}
{{end}}`

func genApi(dir string, api *spec.ApiSpec, isLegacy bool) error {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	err = genApiFile(dir, isLegacy)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dir+strings.ToLower(api.Service.Name+".dart"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()
	t := template.New("apiTemplate")
	t = t.Funcs(funcMap)
	tpl := apiTemplateV2
	if isLegacy {
		tpl = apiTemplate
	}
	t, err = t.Parse(tpl)
	if err != nil {
		return err
	}

	return t.Execute(file, api)
}

func genApiFile(dir string, isLegacy bool) error {
	path := dir + "api.dart"
	if fileExists(path) {
		return nil
	}
	apiFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer apiFile.Close()
	tpl := apiFileContentV2
	if isLegacy {
		tpl = apiFileContent
	}
	_, err = apiFile.WriteString(tpl)
	return err
}
