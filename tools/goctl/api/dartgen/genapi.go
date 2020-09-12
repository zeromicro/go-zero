package dartgen

import (
	"os"
	"text/template"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const apiTemplate = `import 'api.dart';
import '../data/{{with .Info}}{{.Title}}{{end}}.dart';
{{with .Service}}
/// {{.Name}}
{{range .Routes}}
/// --{{.Path}}--
///
/// 请求: {{with .RequestType}}{{.Name}}{{end}}
/// 返回: {{with .ResponseType}}{{.Name}}{{end}}
Future {{pathToFuncName .Path}}( {{if ne .Method "get"}}{{with .RequestType}}{{.Name}} request,{{end}}{{end}}
    {Function({{with .ResponseType}}{{.Name}}{{end}}) ok,
    Function(String) fail,
    Function eventually}) async {
  await api{{if eq .Method "get"}}Get{{else}}Post{{end}}('{{.Path}}',{{if ne .Method "get"}}request,{{end}}
  	 ok: (data) {
    if (ok != null) ok({{with .ResponseType}}{{.Name}}{{end}}.fromJson(data));
  }, fail: fail, eventually: eventually);
}
{{end}}
{{end}}`

func genApi(dir string, api *spec.ApiSpec) error {
	e := os.MkdirAll(dir, 0755)
	if e != nil {
		logx.Error(e)
		return e
	}
	e = genApiFile(dir)
	if e != nil {
		logx.Error(e)
		return e
	}

	file, e := os.OpenFile(dir+api.Info.Title+".dart", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		logx.Error(e)
		return e
	}
	defer file.Close()

	t := template.New("apiTemplate")
	t = t.Funcs(funcMap)
	t, e = t.Parse(apiTemplate)
	if e != nil {
		logx.Error(e)
		return e
	}
	t.Execute(file, api)
	return nil
}

func genApiFile(dir string) error {
	path := dir + "api.dart"
	if fileExists(path) {
		return nil
	}
	apiFile, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		logx.Error(e)
		return e
	}
	defer apiFile.Close()
	apiFile.WriteString(apiFileContent)
	return nil
}
