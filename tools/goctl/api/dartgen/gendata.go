package dartgen

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const dataTemplate = `// --{{with .APISpec.Info}}{{.Title}}{{end}}--
{{ range .APISpec.Types}}
class {{.Name}}{
	{{range .Members}}
	/// {{.Comment}}
	final {{if isNumberType .Type.Name}}num{{else}}{{.Type.Name}}{{end}} {{lowCamelCase .Name}};
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

	{{ range $.InnerClassList}}
	{{.}}
	{{end}}
}
{{end}}
`

const dataTemplateV2 = `// --{{with .APISpec.Info}}{{.Title}}{{end}}--
{{ range .APISpec.Types}}
class {{.Name}} {
	{{range .Members}}
	{{if .Comment}}{{.Comment}}{{end}}
	final {{if isNumberType .Type.Name}}num{{else}}{{.Type.Name}}{{end}} {{lowCamelCase .Name}};
  {{end}}{{.Name}}({{if .Members}}{
	{{range .Members}}  required this.{{lowCamelCase .Name}},
	{{end}}}{{end}});
	factory {{.Name}}.fromJson(Map<String,dynamic> m) {
		return {{.Name}}(
			{{range .Members}}
				{{lowCamelCase .Name}}: {{appendNullCoalescing .}}
					{{if isAtomicType .Type.Name}}
						m['{{getPropertyFromMember .}}'] {{appendDefaultEmptyValue .Type.Name}}
					{{else if isAtomicListType .Type.Name}}
						m['{{getPropertyFromMember .}}']?.cast<{{getCoreType .Type.Name}}>() {{appendDefaultEmptyValue .Type.Name}}
					{{else if isClassListType .Type.Name}}
						((m['{{getPropertyFromMember .}}'] {{appendDefaultEmptyValue .Type.Name}}) as List<dynamic>).map((i) => {{getCoreType .Type.Name}}.fromJson(i)).toList()
					{{else if isMapType .Type.Name}}
						{{if isNumberType .Type.Name}}num{{else}}{{.Type.Name}}{{end}}.from(m['{{getPropertyFromMember .}}'] ?? {})
					{{else}}
						{{.Type.Name}}.fromJson(m['{{getPropertyFromMember .}}']){{end}}
			,{{end}}
		);
	}
	Map<String,dynamic> toJson() {
		return { {{range .Members}}
			'{{getPropertyFromMember .}}': 
				{{if isDirectType .Type.Name}}
					{{lowCamelCase .Name}}
				{{else if isMapType .Type.Name}}
					{{lowCamelCase .Name}}
				{{else if isClassListType .Type.Name}}
					{{lowCamelCase .Name}}{{if isNullableType .Type.Name}}?{{end}}.map((i) => i{{if isListItemsNullable .Type.Name}}?{{end}}.toJson())
				{{else}}
					{{lowCamelCase .Name}}{{if isNullableType .Type.Name}}?{{end}}.toJson()
				{{end}}
			,{{end}}
		};
	}

	{{ range $.InnerClassList}}
	{{.}}
	{{end}}
}
{{end}}`

type DartSpec struct {
	APISpec        *spec.ApiSpec
	InnerClassList []string
}

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

	err, dartSpec := convertDataType(api, isLegacy)
	if err != nil {
		return err
	}

	return t.Execute(file, dartSpec)
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

func convertDataType(api *spec.ApiSpec, isLegacy bool) (error, *DartSpec) {
	var result DartSpec
	types := api.Types
	if len(types) == 0 {
		return nil, &result
	}

	for _, ty := range types {
		defineStruct, ok := ty.(spec.DefineStruct)
		if ok {
			for index, member := range defineStruct.Members {
				structMember, ok := member.Type.(spec.NestedStruct)
				if ok {
					defineStruct.Members[index].Type = spec.PrimitiveType{RawName: member.Name}
					t := template.New("dataTemplate")
					t = t.Funcs(funcMap)
					tpl := dataTemplateV2
					if isLegacy {
						tpl = dataTemplate
					}
					t, err := t.Parse(tpl)
					if err != nil {
						return err, nil
					}

					var innerClassSpec = &spec.ApiSpec{
						Types: []spec.Type{
							spec.DefineStruct{
								RawName: member.Name,
								Members: structMember.Members,
							},
						},
					}
					err, dartSpec := convertDataType(innerClassSpec, isLegacy)
					if err != nil {
						return err, nil
					}

					writer := bytes.NewBuffer(nil)
					err = t.Execute(writer, dartSpec)
					if err != nil {
						return err, nil
					}
					result.InnerClassList = append(result.InnerClassList, writer.String())
				} else {
					tp, err := specTypeToDart(member.Type)
					if err != nil {
						return err, nil
					}
					defineStruct.Members[index].Type = buildSpecType(member.Type, tp)
				}
			}
		}
	}
	result.APISpec = api

	return nil, &result
}
