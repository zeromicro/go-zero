package gen

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/util"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func GenBase(dir string, api *spec.ApiSpec) error {
	data := &template.PyBaseTemplateData{
		ClientName: util.PascalCase(api.Service.Name),
	}
	return template.GenFile(dir, "base.py", template.ApiBase, data)
}
