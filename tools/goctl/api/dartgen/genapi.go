package dartgen

import (
	"os"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
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
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	err = genApiFile(dir)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dir+api.Service.Name+".dart", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer file.Close()
	t := template.New("apiTemplate")
	t = t.Funcs(funcMap)
	t, err = t.Parse(apiTemplate)
	if err != nil {
		return err
	}

	return t.Execute(file, api)
}

func genApiFile(dir string) error {
	path := dir + "api.dart"
	if fileExists(path) {
		return nil
	}
	apiFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer apiFile.Close()
	_, err = apiFile.WriteString(apiFileContent)
	return err
}
