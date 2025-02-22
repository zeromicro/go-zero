package gen

import (
	"regexp"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/util"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func GenActions(api *spec.ApiSpec) ([]*template.ApiActionTemplateData, error) {
	actions := []*template.ApiActionTemplateData{}
	handlerExp := regexp.MustCompile("_handler$")
	for _, g := range api.Service.Groups {
		prefix := g.GetAnnotation("prefix")
		for _, r := range g.Routes {
			name := util.SnakeCase(r.Handler)
			if ok := handlerExp.MatchString(name); ok {
				name = name[0 : len(name)-8]
			}
			a := template.ApiActionTemplateData{
				ActionName: name,
				HttpMethod: strings.ToUpper(r.Method),
				UrlPrefix:  prefix,
				UrlPath:    r.Path,
			}
			if r.RequestType != nil {
				rm, err := GenMessage(r.RequestType)
				if err != nil {
					return nil, err
				}
				a.RequestMessage = rm
			}
			if r.ResponseType != nil {
				rm, err := GenMessage(r.ResponseType)
				if err != nil {
					return nil, err
				}
				a.ResponseMessage = rm
			}
			actions = append(actions, &a)
		}
	}
	return actions, nil
}

func GenClient(dir string, api *spec.ApiSpec) error {
	actions, err := GenActions(api)
	if err != nil {
		return err
	}
	data := template.ApiClientTemplateData{
		ClientName: util.SnakeCase(api.Service.Name),
		Actions:    actions,
	}

	if err := template.GenFile(dir, "base.h", template.ApiBaseHeader, data); err != nil {
		return err
	}

	if err := template.GenFile(dir, "base.c", template.ApiBaseSource, data); err != nil {
		return err
	}

	if err := template.GenFile(dir, "client.h", template.ApiClientHeader, data); err != nil {
		return err
	}

	if err := template.GenFile(dir, "client.c", template.ApiClientSource, data); err != nil {
		return err
	}

	return nil
}
