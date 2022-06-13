package {{.pkgName}}

import (
	{{.imports}}
)

func {{.function}}(ctx context.Context, cc *svc.ClientContext, {{.request}}) {{.responseType}} {
	resp, err := cc.Do(ctx, {{.method}}, fmt.Sprintf("%s%s", cc.Host(), "{{.route}}"), {{.httpRequest}})
	if err != nil {
		{{.returnErrString}}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	{{.returnString}}
}
