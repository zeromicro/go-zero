package {{.pkgName}}

import (
	{{.imports}}
)

func {{.function}}(ctx context.Context, cli *http.Client, host string, {{.request}}) {{.responseType}} {
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}
	svc := httpc.NewServiceWithClient("{{.function}}", cli)
	resp, err := svc.Do(ctx, "{{.method}}", fmt.Sprintf("%s%s", host, "{{.route}}"), req)
	if err != nil {
		{{.returnErrString}}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	{{.returnString}}
}