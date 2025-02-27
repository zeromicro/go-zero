package gen

import (
	"regexp"

	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/util"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func GenClient(dir string, api *spec.ApiSpec) error {
	data := &template.PyClientTemplateData{
		ClientName: util.PascalCase(api.Service.Name),
		Actions:    []*template.PyActionTemplateData{},
	}
	handlerExp := regexp.MustCompile("_handler$")
	for _, g := range api.Service.Groups {
		prefix := g.GetAnnotation("prefix")
		for _, r := range g.Routes {
			name := util.SnakeCase(r.Handler)
			if ok := handlerExp.MatchString(name); ok {
				name = name[0 : len(name)-8]
			}
			a := &template.PyActionTemplateData{
				ActionName: name,
				HttpMethod: r.Method,
				UrlPrefix:  prefix,
				UrlPath:    r.Path,
			}
			if r.RequestType != nil {
				rm, err := GenMessage(r.RequestType)
				if err != nil {
					return err
				}
				a.RequestMessage = rm
			}
			if r.ResponseType != nil {
				rm, err := GenMessage(r.ResponseType)
				if err != nil {
					return err
				}
				a.ResponseMessage = rm
			}
			data.Actions = append(data.Actions, a)
		}
	}

	return template.GenFile(dir, "__init__.py", template.ApiClient, data)
}
