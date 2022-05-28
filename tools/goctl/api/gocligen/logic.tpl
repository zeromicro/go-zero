package {{.pkgName}}

import (
	{{.imports}}
)

func {{.function}}(ctx context.Context, cli *http.Client, host string, {{.request}}) {{.responseType}} {
    u, err := url.Parse(host)
	if err != nil {
		{{.returnErrString}}
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	host = u.String()

	svc := httpc.NewServiceWithClient("{{.function}}", cli)
	resp, err := svc.Do(ctx, "{{.method}}", fmt.Sprintf("%s%s", host, "{{.route}}"), {{.httpRequest}})
	if err != nil {
		{{.returnErrString}}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	{{.returnString}}
}
