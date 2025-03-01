package gen

import "github.com/zeromicro/go-zero/tools/goctl/api/tsgen/template"

func GenRequests(dir string, caller string) error {
	data := template.RequestTemplateData{
		Caller: caller,
	}
	return template.GenTsFile(dir, "gocliRequest", template.Requests, data)
}
