package dartgen

import (
	"os"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const dataTemplate = `// --{{with .Info}}{{.Title}}{{end}}--
{{ range .Types}}
class {{.Name}}{
	{{range .Members}}
	/// {{.Comment}}
	final {{.Type.Name}} {{lowCamelCase .Name}};
	{{end}}
	{{.Name}}({ {{range .Members}}
		this.{{lowCamelCase .Name}},{{end}}
	});
	factory {{.Name}}.fromJson(Map<String,dynamic> m) {
		return {{.Name}}({{range .Members}}
			{{lowCamelCase .Name}}: {{if isDirectType .Type.Name}}m['{{getPropertyFromMember .}}']{{else if isClassListType .Type.Name}}(m['{{getPropertyFromMember .}}'] as List<dynamic>).map((i) => {{getCoreType .Type.Name}}.fromJson(i)){{else}}{{.Type.Name}}.fromJson(m['{{getPropertyFromMember .}}']){{end}},{{end}}
		);
	}
	Map<String,dynamic> toJson() {
		return { {{range .Members}}
			'{{getPropertyFromMember .}}': {{if isDirectType .Type.Name}}{{lowCamelCase .Name}}{{else if isClassListType .Type.Name}}{{lowCamelCase .Name}}.map((i) => i.toJson()){{else}}{{lowCamelCase .Name}}.toJson(){{end}},{{end}}
		};
	}
}
{{end}}
`

const dataTemplateV2 = `// --{{with .Info}}{{.Title}}{{end}}--
{{ range .Types}}
class {{.Name}} {
	{{range .Members}}
	{{if .Comment}}{{.Comment}}{{end}}
	final {{.Type.Name}} {{lowCamelCase .Name}};
  {{end}}{{.Name}}({{if .Members}}{
	{{range .Members}}  required this.{{lowCamelCase .Name}},
	{{end}}}{{end}});
	factory {{.Name}}.fromJson(Map<String,dynamic> m) {
		return {{.Name}}({{range .Members}}
			{{lowCamelCase .Name}}: {{if isDirectType .Type.Name}}m['{{getPropertyFromMember .}}']
			{{else if isClassListType .Type.Name}}(m['{{getPropertyFromMember .}}'] as List<dynamic>).map((i) => {{getCoreType .Type.Name}}.fromJson(i)).toList()
			{{else}}{{.Type.Name}}.fromJson(m['{{getPropertyFromMember .}}']){{end}},{{end}}
		);
	}
	Map<String,dynamic> toJson() {
		return { {{range .Members}}
			'{{getPropertyFromMember .}}': {{if isDirectType .Type.Name}}{{lowCamelCase .Name}}{{else if isClassListType .Type.Name}}{{lowCamelCase .Name}}.map((i) => i.toJson()){{else}}{{lowCamelCase .Name}}.toJson(){{end}},{{end}}
		};
	}
}
{{end}}`

func genData(dir string, api *spec.ApiSpec, isLegacy bool) error {
	err := os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	err = genTokens(dir, isLegacy)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dir+strings.ToLower(api.Service.Name+".dart"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	t := template.New("dataTemplate")
	t = t.Funcs(funcMap)
	tpl := dataTemplateV2
	if isLegacy {
		tpl = dataTemplate
	}
	t, err = t.Parse(tpl)
	if err != nil {
		return err
	}

	err = convertDataType(api)
	if err != nil {
		return err
	}

	return t.Execute(file, api)
}

func genTokens(dir string, isLeagcy bool) error {
	path := dir + "tokens.dart"
	if fileExists(path) {
		return nil
	}

	tokensFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer tokensFile.Close()
	tpl := tokensFileContentV2
	if isLeagcy {
		tpl = tokensFileContent
	}
	_, err = tokensFile.WriteString(tpl)
	return err
}

func convertDataType(api *spec.ApiSpec) error {
	types := api.Types
	if len(types) == 0 {
		return nil
	}

	for _, ty := range types {
		defineStruct, ok := ty.(spec.DefineStruct)
		if ok {
			for index, member := range defineStruct.Members {
				tp, err := specTypeToDart(member.Type)
				if err != nil {
					return err
				}
				defineStruct.Members[index].Type = buildSpecType(member.Type, tp)
			}
		}
	}

	return nil
}
