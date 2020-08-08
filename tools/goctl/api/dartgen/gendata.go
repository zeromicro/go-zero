package dartgen

import (
	"os"
	"text/template"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const dataTemplate = `// --{{with .Info}}{{.Title}}{{end}}--
{{ range .Types}}
class {{.Name}}{
	{{range .Members}}
	/// {{.Comment}}
	final {{.Type}} {{lowCamelCase .Name}};
	{{end}}
	{{.Name}}({ {{range .Members}}
		this.{{lowCamelCase .Name}},{{end}}
	});
	factory {{.Name}}.fromJson(Map<String,dynamic> m) {
		return {{.Name}}({{range .Members}}
			{{lowCamelCase .Name}}: {{if isDirectType .Type}}m['{{tagGet .Tag "json"}}']{{else if isClassListType .Type}}(m['{{tagGet .Tag "json"}}'] as List<dynamic>).map((i) => {{getCoreType .Type}}.fromJson(i)){{else}}{{.Type}}.fromJson(m['{{tagGet .Tag "json"}}']){{end}},{{end}}
		);
	}
	Map<String,dynamic> toJson() {
		return { {{range .Members}}
			'{{tagGet .Tag "json"}}': {{if isDirectType .Type}}{{lowCamelCase .Name}}{{else if isClassListType .Type}}{{lowCamelCase .Name}}.map((i) => i.toJson()){{else}}{{lowCamelCase .Name}}.toJson(){{end}},{{end}}
		};
	}
}
{{end}}
`

func genData(dir string, api *spec.ApiSpec) error {
	e := os.MkdirAll(dir, 0755)
	if e != nil {
		logx.Error(e)
		return e
	}
	e = genTokens(dir)
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

	t := template.New("dataTemplate")
	t = t.Funcs(funcMap)
	t, e = t.Parse(dataTemplate)
	if e != nil {
		logx.Error(e)
		return e
	}

	convertMemberType(api)
	return t.Execute(file, api)
}

func genTokens(dir string) error {
	path := dir + "tokens.dart"
	if fileExists(path) {
		return nil
	}
	tokensFile, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		logx.Error(e)
		return e
	}
	defer tokensFile.Close()
	tokensFile.WriteString(tokensFileContent)
	return nil
}
