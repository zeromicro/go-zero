package phpgen

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

//go:embed ApiBaseClient.php
var apiBaseClientTemplate string

//go:embed ApiException.php
var apiExceptionTemplate string

func genClient(dir string, ns string, api *spec.ApiSpec) error {
	if err := writeBaseClient(dir, ns); err != nil {
		return err
	}
	if err := writeException(dir, ns); err != nil {
		return err
	}
	return writeClient(dir, ns, api)
}

func writeBaseClient(dir string, ns string) error {
	bPath := filepath.Join(dir, "ApiBaseClient.php")
	bHead := fmt.Sprintf("<?php\n\nnamespace %s;", ns)
	bSrc := strings.Replace(apiBaseClientTemplate, "<?php", bHead, 1)

	return os.WriteFile(bPath, []byte(bSrc), 0644)
}

func writeException(dir string, ns string) error {
	ePath := filepath.Join(dir, "ApiException.php")
	eHead := fmt.Sprintf("<?php\n\nnamespace %s;", ns)
	eSrc := strings.Replace(apiExceptionTemplate, "<?php", eHead, 1)

	return os.WriteFile(ePath, []byte(eSrc), 0644)
}

func writeClient(dir string, ns string, api *spec.ApiSpec) error {
	name := camelCase(api.Service.Name, true)
	fp := filepath.Join(dir, fmt.Sprintf("%sClient.php", name))
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	// 头部
	fmt.Fprintf(f, "<?php \n\nnamespace %s;\n\n", ns)

	// 类

	fmt.Fprintf(f, "class %sClient extends ApiBaseClient {\n", name)

	for _, g := range api.Service.Groups {
		prefix := g.GetAnnotation("prefix")
		p := camelCase(prefix, true)

		// 路由
		for _, r := range g.Routes {
			an := camelCase(r.Path, true)

			writeIndent(f, 4)
			fmt.Fprintf(f, "public function %s%s%s(", strings.ToLower(r.Method), p, an)
			if r.RequestType != nil {
				fmt.Fprint(f, "$request")
			}
			fmt.Fprintln(f, ") {")

			writeIndent(f, 8)
			fmt.Fprintf(f, "$result = $this->request('%s%s', '%s',", prefix, r.Path, strings.ToLower(r.Method))

			if r.RequestType != nil {
				params := []string{}
				for _, tagKey := range tagKeys {
					if hasTagMembers(r.RequestType, tagKey) {
						sn := camelCase(fmt.Sprintf("get-%s", tagToSubName(tagKey)), false)
						params = append(params, fmt.Sprintf("$request->%s()", sn))
					} else {
						params = append(params, "null")
					}
				}
				fmt.Fprint(f, strings.Join(params, ","))
			} else {
				fmt.Fprint(f, "null, null, null, null")
			}

			fmt.Fprintln(f, ");")

			writeIndent(f, 8)
			if r.ResponseType != nil {
				n := camelCase(r.ResponseType.Name(), true)
				fmt.Fprintf(f, "$response = new %s();\n", n)
				definedType, ok := r.ResponseType.(spec.DefineStruct)
				if !ok {
					return fmt.Errorf("type %s not supported", n)
				}
				if err := writeResponseHeader(f, &definedType); err != nil {
					return err
				}
				if err := writeResponseBody(f, &definedType); err != nil {
					return err
				}
				writeIndent(f, 8)
				fmt.Fprint(f, "return $response;\n")
			} else {
				fmt.Fprint(f, "return null;\n")
			}

			writeIndent(f, 4)
			fmt.Fprintln(f, "}")
		}
	}

	fmt.Fprintln(f, "}")

	return nil
}

func writeResponseBody(f *os.File, definedType *spec.DefineStruct) error {
	// 获取字段
	ms := definedType.GetTagMembers(bodyTagKey)
	if len(ms) <= 0 {
		return nil
	}
	writeIndent(f, 8)
	fmt.Fprint(f, "$response->getBody()")
	for _, m := range ms {
		tags := m.Tags()
		k := ""
		if len(tags) > 0 {
			k = tags[0].Name
		} else {
			k = m.Name
		}
		fmt.Fprintf(f, "\n            ->set%s($result['body']['%s'])", camelCase(m.Name, true), k)
	}
	fmt.Fprintln(f, ";")
	return nil
}

func writeResponseHeader(f *os.File, definedType *spec.DefineStruct) error {
	// 获取字段
	ms := definedType.GetTagMembers(headerTagKey)
	if len(ms) <= 0 {
		return nil
	}
	writeIndent(f, 8)
	fmt.Fprint(f, "$response->getHeader()")
	for _, m := range ms {
		tags := m.Tags()
		k := ""
		if len(tags) > 0 {
			k = tags[0].Name
		} else {
			k = m.Name
		}
		fmt.Fprintf(f, "\n            ->set%s($result['header']['%s'])", camelCase(m.Name, true), k)
	}
	fmt.Fprintln(f, ";")
	return nil
}
