package unigen

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

//go:embed ApiBaseClient.ts
var apiBaseClientTemplate string

func writeTemplate(dir string, name string, template string) error {
	p := filepath.Join(dir, fmt.Sprintf("%s.ts", name))
	return os.WriteFile(p, []byte(template), 0644)
}

func genClient(dir string, api *spec.ApiSpec) error {
	if err := writeTemplate(dir, "ApiBaseClient", apiBaseClientTemplate); err != nil {
		return err
	}

	return writeClient(dir, api)
}

func writeClient(dir string, api *spec.ApiSpec) error {
	name := camelCase(api.Service.Name, true)
	fp := filepath.Join(dir, fmt.Sprintf("%sClient.ts", name))
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	// 引入
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			if r.RequestType != nil {
				rn := r.RequestType.Name()
				fmt.Fprintf(f, "import type { %s } from './%s';\n", rn, rn)
			}

			if r.ResponseType != nil {
				rn := r.ResponseType.Name()
				fmt.Fprintf(f, "import type { %s", rn)
				for _, tagKey := range tagKeys {
					if hasTagMembers(r.ResponseType, tagKey) {
						sn := camelCase(fmt.Sprintf("%s-%s", rn, tagToSubName(tagKey)), true)
						fmt.Fprintf(f, ", %s", sn)
					}
				}
				fmt.Fprintf(f, "} from './%s';\n", rn)
			}
		}
	}
	fmt.Fprintf(f, "import { ApiBaseClient } from './ApiBaseClient';\n\n")

	// 类
	fmt.Fprintf(f, "export class %sClient extends ApiBaseClient {\n", name)

	// 方法
	for _, g := range api.Service.Groups {
		prefix := g.GetAnnotation("prefix")
		p := camelCase(prefix, true)

		// 路由
		for _, r := range g.Routes {
			an := camelCase(r.Path, true)
			method := strings.ToLower(r.Method)

			writeIndent(f, 4)
			fmt.Fprintf(f, "async %s%s%s(", method, p, an)

			if r.RequestType != nil {
				fmt.Fprintf(f, "request: %s, body?: any", r.RequestType.Name())
			} else {
				fmt.Fprintf(f, "body?: any")
			}

			if r.ResponseType != nil {
				fmt.Fprintf(f, "): Promise<%s> {\n", r.ResponseType.Name())
			} else {
				fmt.Fprintf(f, "): Promise<UniApp.RequestSuccessCallbackResult> {\n")
			}

			writeIndent(f, 8)
			fmt.Fprintf(f, "const response = await this.request('%s', '%s%s',", strings.ToUpper(method), prefix, r.Path)
			if hasTagMembers(r.RequestType, formTagKey) {
				fmt.Fprint(f, " request.query,")
			} else {
				fmt.Fprint(f, " undefined,")
			}
			if hasTagMembers(r.RequestType, headerTagKey) {
				fmt.Fprint(f, " request.header,")
			} else {
				fmt.Fprint(f, " undefined,")
			}
			if hasTagMembers(r.RequestType, bodyTagKey) {
				fmt.Fprint(f, " body ?? request.body")
			} else {
				fmt.Fprint(f, " body")
			}
			fmt.Fprint(f, ");\n")

			if r.ResponseType != nil {
				writeIndent(f, 8)
				fmt.Fprintf(f, "const result: %s = {\n", r.ResponseType.Name())
				writeIndent(f, 12)
				if hasTagMembers(r.ResponseType, bodyTagKey) {
					sn := camelCase(fmt.Sprintf("%s-%s", r.ResponseType.Name(), tagToSubName(bodyTagKey)), true)
					fmt.Fprintf(f, "body: response.data as %s\n", sn)
				}
				if hasTagMembers(r.ResponseType, headerTagKey) {
					sn := camelCase(fmt.Sprintf("%s-%s", r.ResponseType.Name(), tagToSubName(headerTagKey)), true)
					fmt.Fprintf(f, "header: response.header as %s\n", sn)
				}
				writeIndent(f, 8)
				fmt.Fprint(f, "};\n")
				writeIndent(f, 8)
				fmt.Fprintf(f, "return result;\n")
			} else {
				writeIndent(f, 8)
				fmt.Fprintf(f, "return response;\n")
			}

			writeIndent(f, 4)
			fmt.Fprintf(f, "}\n\n")
		}
	}

	fmt.Fprintf(f, "}\n")

	return nil
}
