package phpgen

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const (
	formTagKey   = "form"
	pathTagKey   = "path"
	headerTagKey = "header"
	bodyTagKey   = "json"
)

var (
	// 这个顺序与 PHP ApiBaseClient request 参数相关，改动时要注意 PHP 那边的代码。
	tagKeys = []string{pathTagKey, formTagKey, headerTagKey, bodyTagKey}
)

func tagToSubName(tagKey string) string {
	suffix := tagKey
	switch tagKey {
	case "json":
		suffix = "body"
	case "form":
		suffix = "query"
	}
	return suffix
}

func getMessageName(tn string, tagKey string, isPascal bool) string {
	suffix := tagToSubName(tagKey)
	return camelCase(fmt.Sprintf("%s-%s", tn, suffix), isPascal)
}

func hasTagMembers(t spec.Type, tagKey string) bool {
	definedType, ok := t.(spec.DefineStruct)
	if !ok {
		return false
	}
	ms := definedType.GetTagMembers(tagKey)
	return len(ms) > 0
}

func genMessages(dir string, ns string, api *spec.ApiSpec) error {
	for _, t := range api.Types {
		tn := t.Name()
		definedType, ok := t.(spec.DefineStruct)
		if !ok {
			return fmt.Errorf("type %s not supported", tn)
		}

		// 子类型
		tags := []string{}
		for _, tagKey := range tagKeys {
			// 获取字段
			ms := definedType.GetTagMembers(tagKey)
			if len(ms) <= 0 {
				continue
			}

			// 打开文件
			cn := getMessageName(tn, tagKey, true)
			tags = append(tags, tagKey)
			fp := filepath.Join(dir, fmt.Sprintf("%s.php", cn))
			f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return err
			}
			defer f.Close()

			// 写入
			if err := writeSubMessage(f, ns, cn, ms); err != nil {
				return err
			}
		}

		// 主类型
		rn := camelCase(tn, true)
		fp := filepath.Join(dir, fmt.Sprintf("%s.php", rn))
		f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := writeMessage(f, ns, rn, tags); err != nil {
			return nil
		}
	}

	return nil
}

func writeMessage(f *os.File, ns string, rn string, tags []string) error {
	// 文件头
	fmt.Fprintln(f, "<?php")

	// 名字空间
	fmt.Fprintf(f, "namespace %s;\n\n", ns)

	// 类
	fmt.Fprintf(f, "class %s {\n", rn)

	// 字段
	for _, tag := range tags {
		writeIndent(f, 4)
		fmt.Fprintf(f, "private $%s;\n", tagToSubName(tag))
	}

	// 构造函数
	writeIndent(f, 4)
	fmt.Fprint(f, "public function __construct(){\n")
	for _, tag := range tags {
		sn := tagToSubName(tag)
		cn := getMessageName(rn, tag, true)
		writeIndent(f, 8)
		fmt.Fprintf(f, "$this->%s = new %s();\n", sn, cn)
	}
	writeIndent(f, 4)
	fmt.Fprintln(f, "}")

	// getter
	for _, tag := range tags {
		sn := tagToSubName(tag)
		pn := camelCase(fmt.Sprintf("get-%s", sn), false)
		writeIndent(f, 4)
		fmt.Fprintf(f, "public function %s() { return $this->%s; }\n", pn, sn)
	}

	fmt.Fprintln(f, "}")

	return nil
}

func writeSubMessage(f *os.File, ns string, cn string, ms []spec.Member) error {
	// 文件头
	fmt.Fprintln(f, "<?php")

	// 名字空间
	fmt.Fprintf(f, "namespace %s;\n\n", ns)

	// 类
	fmt.Fprintf(f, "class %s {\n", cn)

	// 字段
	for _, m := range ms {
		writeField(f, m)
	}

	// getter setter
	for _, m := range ms {
		writeProperty(f, m)
	}

	// toQueryString
	writeIndent(f, 4)
	fmt.Fprintf(f, "public function toQueryString() {\n        return http_build_query([\n")
	writeMembersToPhpArrayItems(f, ms, 12)
	fmt.Fprintln(f, "        ]);\n    }")

	// toJsonString
	writeIndent(f, 4)
	fmt.Fprintf(f, "public function toJsonString() {\n        return json_encode([\n")
	writeMembersToPhpArrayItems(f, ms, 12)
	fmt.Fprintln(f, "        ], JSON_UNESCAPED_UNICODE);\n    }")

	// toAssocArray
	writeIndent(f, 4)
	fmt.Fprintf(f, "public function toAssocArray() {\n        return [\n")
	writeMembersToPhpArrayItems(f, ms, 12)
	fmt.Fprintln(f, "        ];\n    }")

	_, err := fmt.Fprintln(f, "}")
	return err
}

func writeMembersToPhpArrayItems(f *os.File, ms []spec.Member, indent int) {
	for _, m := range ms {
		tags := m.Tags()
		n := camelCase(m.Name, false)
		k := ""
		if len(tags) > 0 {
			k = tags[0].Name
		} else {
			k = n
		}
		writeIndent(f, indent)
		fmt.Fprintf(f, "'%s' => $this->%s,\n", k, n)
	}
}

func writeField(f *os.File, m spec.Member) {
	writeIndent(f, 4)
	fmt.Fprintf(f, "private $%s;\n", camelCase(m.Name, false))
}

func writeProperty(f *os.File, m spec.Member) {
	pName := camelCase(m.Name, true)
	cName := camelCase(m.Name, false)
	writeIndent(f, 4)
	fmt.Fprintf(f, "public function get%s() { return $this->%s; }\n\n", pName, cName)
	writeIndent(f, 4)
	fmt.Fprintf(f, "public function set%s($v) { $this->%s = $v; return $this; }\n\n", pName, cName)
}

func writeIndent(f *os.File, n int) {
	for i := 0; i < n; i++ {
		fmt.Fprint(f, " ")
	}
}
