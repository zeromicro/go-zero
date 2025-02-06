package csgen

import (
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/api/csgen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/csgen/util"
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

func genMessages(dir string, ns string, api *spec.ApiSpec) error {
	for _, t := range api.Types {
		tn := t.Name()
		definedType, ok := t.(spec.DefineStruct)
		if !ok {
			return fmt.Errorf("type %s not supported", tn)
		}

		cn := util.CamelCase(tn, true)

		data := template.CSharpApiMessageTemplateData{
			CSharpTemplateData: template.CSharpTemplateData{Namespace: ns},
			MessageName:        cn,
			Fields:             []template.CSharpApiMessageFieldTemplateData{},
		}

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
				tn, err := apiTypeToCsTypeName(m.Type)
				if err != nil {
					return err
				}

				f := template.CSharpApiMessageFieldTemplateData{
					FieldName:  util.CamelCase(m.Name, true),
					KeyName:    k,
					TypeName:   tn,
					IsOptional: m.IsOptionalOrOmitEmpty(),
					Tag:        tagKey,
				}
				data.Fields = append(data.Fields, f)
			}
		}

		if err := template.WriteFile(dir, cn, template.ApiMessage, data); err != nil {
			return err
		}
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
	case spec.DefineStruct:
		return t.Name(), nil
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
		return util.UpperHead(tp, 1), true
	case "uint16", "uint32", "uint64":
		return util.UpperHead(tp, 2), true
	case "float", "float32":
		return "float", true
	case "float64":
		return "double", true
	}
	return "", false
}
