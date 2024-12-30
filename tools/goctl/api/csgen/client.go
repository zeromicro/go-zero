package csgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

//go:embed ApiAttribute.cs
var apiAttributeTemplate string

//go:embed ApiBaseClient.cs
var apiApiBaseClientTemplate string

func genClient(dir string, ns string, api *spec.ApiSpec) error {
	if err := writeTemplate(dir, ns, "ApiAttribute", apiAttributeTemplate); err != nil {
		return err
	}
	if err := writeTemplate(dir, ns, "ApiBaseClient", apiApiBaseClientTemplate); err != nil {
		return err
	}

	return writeClient(dir, ns, api)
}

func writeTemplate(dir string, ns string, name string, template string) error {
	fp := filepath.Join(dir, fmt.Sprintf("%s.cs", name))
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "namespace %s;\r\n\r\n", ns)
	fmt.Fprint(f, template)
	return nil
}

func writeClient(dir string, ns string, api *spec.ApiSpec) error {
	name := camelCase(api.Service.Name, true)
	fp := filepath.Join(dir, fmt.Sprintf("%sClient.cs", name))
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "namespace %s;\r\n\r\n", ns)

	// 类
	fmt.Fprintf(f, "public sealed class %sClient : ApiBaseClient\r\n{\r\n", name)

	// 构造函数
	fmt.Fprintf(f, "    public %sClient(string host, short port, string scheme = \"http\") : base(host, port, scheme){}\r\n", name)

	// 组
	for _, g := range api.Service.Groups {
		prefix := g.GetAnnotation("prefix")
		p := camelCase(prefix, true)

		// 路由
		for _, r := range g.Routes {
			an := camelCase(r.Path, true)
			method := upperHead(strings.ToLower(r.Method), 1)

			writeIndent(f, 4)
			fmt.Fprint(f, "public async ")
			if r.ResponseType != nil {
				fmt.Fprintf(f, "Task<%s>", r.ResponseType.Name())
			} else {
				fmt.Fprint(f, "Task<HttpResponseMessage>")
			}
			fmt.Fprintf(f, " %s%s%sAsync(", method, p, an)
			if r.RequestType != nil {
				fmt.Fprintf(f, "%s request,", r.RequestType.Name())
			}
			fmt.Fprint(f, "CancellationToken cancellationToken)\r\n    {\r\n")

			writeIndent(f, 8)
			fmt.Fprint(f, "return await ")

			if r.RequestType != nil {
				if r.ResponseType != nil {
					fmt.Fprintf(f, "RequestResultAsync<%s,%s>(HttpMethod.%s, \"%s\", request, cancellationToken);\r\n", r.RequestType.Name(), r.ResponseType.Name(), method, r.Path)
				} else {
					fmt.Fprintf(f, "RequestAsync(HttpMethod.%s, \"%s\", request, cancellationToken);\r\n", method, r.Path)
				}
			} else {
				if r.ResponseType != nil {
					fmt.Fprintf(f, "CallResultAsync<%s>(HttpMethod.%s, \"%s\", cancellationToken);\r\n", r.ResponseType.Name(), method, r.Path)
				} else {
					fmt.Fprintf(f, "CallAsync(HttpMethod.%s, \"%s\", cancellationToken);\r\n", method, r.Path)
				}
			}

			writeIndent(f, 4)
			fmt.Fprint(f, "}\r\n")
		}
	}

	fmt.Fprint(f, "}\r\n")

	return nil
}
