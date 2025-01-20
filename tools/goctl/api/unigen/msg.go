package unigen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const (
	formTagKey   = "form"
	pathTagKey   = "path"
	headerTagKey = "header"
	bodyTagKey   = "json"
)

var (
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

func getMessageName(tn string, tagKey string) string {
	suffix := tagToSubName(tagKey)
	return camelCase(fmt.Sprintf("%s-%s", tn, suffix), true)
}

func hasTagMembers(t spec.Type, tagKey string) bool {
	definedType, ok := t.(spec.DefineStruct)
	if !ok {
		return false
	}
	ms := definedType.GetTagMembers(tagKey)
	return len(ms) > 0
}

func writeSubMessage(f *os.File, cn string, ms []spec.Member) error {
	fmt.Fprintf(f, "export type %s = {\n", cn)

	for _, m := range ms {
		writeIndent(f, 4)
		tags := m.Tags()
		k := ""
		if len(tags) > 0 {
			k = tags[0].Name
		} else {
			k = m.Name
		}

		if strings.Contains(k, "-") {
			k = fmt.Sprintf("\"%s\"", k)
		}

		tn, err := apiTypeToUniTsTypeName(m.Type)
		if err != nil {
			return err
		}
		optionalTag := ""
		if m.IsOptionalOrOmitEmpty() {
			optionalTag = "?"
		}
		fmt.Fprintf(f, "%s%s: %s;\n", k, optionalTag, tn)
	}

	fmt.Fprintf(f, "};\n")
	return nil
}

func writeIndent(f *os.File, n int) {
	for i := 0; i < n; i++ {
		fmt.Fprint(f, " ")
	}
}

func genMessages(dir string, api *spec.ApiSpec) error {
	for _, t := range api.Types {
		tn := t.Name()
		definedType, ok := t.(spec.DefineStruct)
		if !ok {
			return fmt.Errorf("type %s not supported", tn)
		}

		// 主类型
		rn := camelCase(tn, true)
		fp := filepath.Join(dir, fmt.Sprintf("%s.ts", rn))
		f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()

		// 子类型
		tags := []string{}
		for _, tagKey := range tagKeys {
			// 获取字段
			ms := definedType.GetTagMembers(tagKey)
			if len(ms) <= 0 {
				continue
			}
			tags = append(tags, tagKey)
			mn := getMessageName(rn, tagKey)
			writeSubMessage(f, mn, ms)
		}

		fmt.Fprintf(f, "export type %s = {\n", tn)

		// 子字段
		for _, tag := range tags {
			// 获取字段
			sn := tagToSubName(tag)
			mn := getMessageName(rn, tag)
			writeIndent(f, 4)
			fmt.Fprintf(f, "%s: %s;\n", sn, mn)
		}

		fmt.Fprint(f, "};\n")
	}
	return nil
}

func apiTypeToUniTsTypeName(t spec.Type) (string, error) {
	switch tt := t.(type) {
	case spec.PrimitiveType:
		r, ok := primitiveType(t.Name())
		if !ok {
			return "", errors.New("unsupported primitive type " + t.Name())
		}

		return r, nil
	case spec.ArrayType:
		et, err := apiTypeToUniTsTypeName(tt.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Array<%s>", et), nil
	case spec.MapType:
		vt, err := apiTypeToUniTsTypeName(tt.Value)
		if err != nil {
			return "", err
		}
		kt, ok := primitiveType(tt.Key)
		if !ok {
			return "", errors.New("unsupported key is not primitive type " + t.Name())
		}
		return fmt.Sprintf("{ [key: %s]: %s; }", kt, vt), nil
	}

	return "", errors.New("unsupported type " + t.Name())
}

func primitiveType(tp string) (string, bool) {
	switch tp {
	case "string":
		return tp, true
	case "bool":
		return "boolean", true
	case "int8", "int", "uint", "byte", "uint8", "int16", "int32", "int64", "uint16", "uint32", "uint64", "float", "float32", "float64":
		return "number", true
	case "[]byte":
		return "Blob", true
	}

	return "", false
}
