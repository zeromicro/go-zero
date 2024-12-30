package csgen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/util"
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

func genMessages(dir string, ns string, api *spec.ApiSpec) error {
	for _, t := range api.Types {
		tn := t.Name()
		definedType, ok := t.(spec.DefineStruct)
		if !ok {
			return fmt.Errorf("type %s not supported", tn)
		}

		cn := camelCase(tn, true)
		fp := filepath.Join(dir, fmt.Sprintf("%s.cs", cn))
		f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()

		// 引入依赖
		fmt.Fprint(f, "using System.Text.Json.Serialization;\r\n\r\n")

		// 名字空间
		fmt.Fprintf(f, "namespace %s;\r\n\r\n", ns)

		// 类
		fmt.Fprintf(f, "public class %s\r\n{\r\n", cn)

		for _, tagKey := range tagKeys {
			// 获取字段
			ms := definedType.GetTagMembers(tagKey)
			if len(ms) <= 0 {
				continue
			}

			for _, m := range ms {
				tags := m.Tags()
				k := ""
				if len(tags) > 0 {
					k = tags[0].Name
				} else {
					k = m.Name
				}

				writeIndent(f, 4)
				switch tagKey {
				case bodyTagKey:
					fmt.Fprintf(f, "[JsonPropertyName(\"%s\")]\r\n", k)
				case headerTagKey:
					fmt.Fprint(f, "[JsonIgnore]\r\n")
					writeIndent(f, 4)
					fmt.Fprintf(f, "[HeaderPropertyName(\"%s\")]\r\n", k)
				case formTagKey:
					fmt.Fprint(f, "[JsonIgnore]\r\n")
					writeIndent(f, 4)
					fmt.Fprintf(f, "[FormPropertyName(\"%s\")]\r\n", k)
				case pathTagKey:
					fmt.Fprint(f, "[JsonIgnore]\r\n")
					writeIndent(f, 4)
					fmt.Fprintf(f, "[PathPropertyName(\"%s\")]\r\n", k)
				}

				writeIndent(f, 4)
				tn, err := apiTypeToCsTypeName(m.Type)
				if err != nil {
					return err
				}

				optionalTag := ""
				if m.IsOptionalOrOmitEmpty() {
					optionalTag = "?"
				}
				fmt.Fprintf(f, "public %s%s %s { get; set; }\r\n\r\n", tn, optionalTag, camelCase(m.Name, true))
			}
		}

		fmt.Fprint(f, "}\r\n")
	}
	return nil
}

func apiTypeToCsTypeName(t spec.Type) (string, error) {
	switch tt := t.(type) {
	case spec.PrimitiveType:
		r, ok := primitiveType(t.Name())
		if !ok {
			return "", errors.New("unsupported primitive type " + t.Name())
		}

		return r, nil
	case spec.ArrayType:
		et, err := apiTypeToCsTypeName(tt.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("List<%s>", et), nil
	case spec.MapType:
		vt, err := apiTypeToCsTypeName(tt.Value)
		if err != nil {
			return "", err
		}
		kt, ok := primitiveType(tt.Key)
		if !ok {
			return "", errors.New("unsupported key is not primitive type " + t.Name())
		}
		return fmt.Sprintf("Dictionary<%s,%s>", kt, vt), nil
	}

	return "", errors.New("unsupported type " + t.Name())
}

func primitiveType(tp string) (string, bool) {
	switch tp {
	case "string", "int", "uint", "bool", "byte":
		return tp, true
	case "int8":
		return "SByte", true
	case "uint8":
		return "byte", true
	case "int16", "int32", "int64":
		return util.UpperFirst(tp), true
	case "uint16", "uint32", "uint64":
		return upperHead(tp, 2), true
	case "float", "float32":
		return "float", true
	case "float64":
		return "double", true
	}
	return "", false
}
