package unigen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/unigen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/unigen/util"
)

func genClient(dir string, api *spec.ApiSpec) error {
	if err := template.WriteFile(dir, "ApiBaseClient", template.ApiBaseClient, nil); err != nil {
		return err
	}

	return writeClient(dir, api)
}

func writeClient(dir string, api *spec.ApiSpec) error {
	name := util.CamelCase(api.Service.Name, true)

	data := template.UniAppApiClientTemplateData{
		ClientName:       name,
		RequestTypes:     []string{},
		ResponseTypes:    []string{},
		ResponseSubTypes: map[string][]string{},
		Routes:           []template.UniAppApiClientRouteTemplateData{},
	}

	// 组
	for _, g := range api.Service.Groups {
		prefix := g.GetAnnotation("prefix")
		p := util.CamelCase(prefix, true)

		// 路由
		for _, r := range g.Routes {
			an := util.CamelCase(r.Path, true)
			method := strings.ToLower(r.Method)

			route := template.UniAppApiClientRouteTemplateData{
				HttpMethod:   method,
				Prefix:       prefix,
				UrlPath:      r.Path,
				ActionPrefix: p,
				ActionName:   an,
			}

			if r.RequestType != nil {
				rn := r.RequestType.Name()
				route.RequestType = &rn
				data.RequestTypes = append(data.RequestTypes, rn)
			}

			if r.ResponseType != nil {
				rn := r.ResponseType.Name()
				route.ResponseType = &rn
				data.ResponseTypes = append(data.ResponseTypes, rn)
				for _, tagKey := range tagKeys {
					if hasTagMembers(r.ResponseType, tagKey) {
						sn := util.CamelCase(fmt.Sprintf("%s-%s", rn, tagToSubName(tagKey)), true)
						data.ResponseSubTypes[rn] = append(data.ResponseSubTypes[rn], sn)
					}
				}
			}

			route.RequestHasQueryString = hasTagMembers(r.RequestType, formTagKey)
			route.RequestHasHeaders = hasTagMembers(r.RequestType, headerTagKey)
			route.RequestHasBody = hasTagMembers(r.RequestType, bodyTagKey)

			if r.ResponseType != nil {
				if hasTagMembers(r.ResponseType, bodyTagKey) {
					sn := util.CamelCase(fmt.Sprintf("%s-%s", r.ResponseType.Name(), tagToSubName(bodyTagKey)), true)
					route.ResponseBodyType = &sn
				}
				if hasTagMembers(r.ResponseType, headerTagKey) {
					sn := util.CamelCase(fmt.Sprintf("%s-%s", r.ResponseType.Name(), tagToSubName(headerTagKey)), true)
					route.ResponseHeadersType = &sn
				}
			}

			data.Routes = append(data.Routes, route)
		}
	}

	return template.WriteFile(dir, fmt.Sprintf("%sClient", name), template.ApiClient, data)
}
