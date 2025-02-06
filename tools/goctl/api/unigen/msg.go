package unigen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/unigen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/unigen/util"
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
	return util.CamelCase(fmt.Sprintf("%s-%s", tn, suffix), true)
}

func hasTagMembers(t spec.Type, tagKey string) bool {
	definedType, ok := t.(spec.DefineStruct)
	if !ok {
		return false
	}
	ms := definedType.GetTagMembers(tagKey)
	return len(ms) > 0
}

func genMessages(dir string, api *spec.ApiSpec) error {
	for _, t := range api.Types {
		tn := t.Name()
		definedType, ok := t.(spec.DefineStruct)
		if !ok {
			return fmt.Errorf("type %s not supported", tn)
		}

		// 主类型
		rn := util.CamelCase(tn, true)
		data := template.UniAppApiMessageTemplateData{
			MessageName: rn,
			SubMessages: []template.UniAppApiSubMessageTemplateData{},
			ImportTypes: []string{},
		}

		for _, tagKey := range tagKeys {
			// 获取字段
			ms := definedType.GetTagMembers(tagKey)
			if len(ms) <= 0 {
				continue
			}

			// 子类型
			sn := tagToSubName(tagKey)
			mn := getMessageName(rn, tagKey)

			data.Fields = append(data.Fields, template.UniAppApiMessageFieldTemplateData{
				FieldName:  sn,
				TypeName:   mn,
				IsOptional: false,
			})

			subMsg := template.UniAppApiSubMessageTemplateData{
				MessageName: mn,
			}
			for _, m := range ms {
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

				tn, b, err := apiTypeToUniTsTypeName(m.Type)
				if err != nil {
					return err
				}

				if len(b) > 0 {
					data.ImportTypes = append(data.ImportTypes, b...)
				}

				f := template.UniAppApiMessageFieldTemplateData{
					FieldName:  k,
					IsOptional: m.IsOptionalOrOmitEmpty(),
					TypeName:   tn,
				}
				subMsg.Fields = append(subMsg.Fields, f)
			}

			data.SubMessages = append(data.SubMessages, subMsg)
		}

		if err := template.WriteFile(dir, rn, template.ApiMessage, data); err != nil {
			return err
		}
	}
	return nil
}

func apiTypeToUniTsTypeName(t spec.Type) (string, []string, error) {
	switch tt := t.(type) {
	case spec.PrimitiveType:
		r, ok := primitiveType(t.Name())
		if !ok {
			return "", []string{}, errors.New("unsupported primitive type " + t.Name())
		}
		return r, []string{}, nil
	case spec.ArrayType:
		et, b, err := apiTypeToUniTsTypeName(tt.Value)
		if err != nil {
			return "", b, err
		}
		return fmt.Sprintf("Array<%s>", et), b, nil
	case spec.MapType:
		vt, b, err := apiTypeToUniTsTypeName(tt.Value)
		if err != nil {
			return "", b, err
		}
		kt, ok := primitiveType(tt.Key)
		if !ok {
			return "", b, errors.New("unsupported key is not primitive type " + t.Name())
		}
		return fmt.Sprintf("{ [key: %s]: %s; }", kt, vt), b, nil
	case spec.DefineStruct:
		return t.Name(), []string{t.Name()}, nil
	}

	return "", []string{}, errors.New("unsupported type " + t.Name())
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
