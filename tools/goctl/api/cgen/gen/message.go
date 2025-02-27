package gen

import (
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/util"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var requestFieldTags = []string{"header", "json", "path", "form"}

func primitiveType(t string) (string, error) {
	switch t {
	case "int", "bool":
		return t, nil
	case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
		return fmt.Sprintf("%s_t", t), nil
	case "string":
		return "char*", nil
	}
	return "", errors.New("no supported of primitive type " + t)
}

func structType(t string) string {
	return fmt.Sprintf("struct __%s_t", util.SnakeCase(t))
}

func formatPrimitiveTag(t string) (string, error) {
	switch t {
	case "bool":
		return "%d", nil
	case "int", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
		return "%d", nil
	case "string":
		return "%s", nil
	}
	return "", errors.New("no supported of primitive type " + t)
}

func GenFormatTag(t spec.Type) (string, error) {
	switch t.(type) {
	case spec.PrimitiveType:
		return formatPrimitiveTag(t.Name())
	case spec.DefineStruct:
		return "", nil
	}
	return "", errors.New("no supported of format type " + t.Name())
}

func GenType(t spec.Type) (string, error) {
	switch t.(type) {
	case spec.PrimitiveType:
		return primitiveType(t.Name())
	case spec.ArrayType:
		return "array_t", nil
	case spec.DefineStruct:
		return structType(t.Name()), nil
	}
	return "", errors.New("no supported of type " + t.Name())
}

func GenFieldsByTag(definedType spec.DefineStruct, tag string) ([]*template.ApiFieldTemplateData, error) {
	result := []*template.ApiFieldTemplateData{}
	for _, m := range definedType.GetTagMembers(tag) {
		ft, err := GenType(m.Type)
		if err != nil {
			return nil, err
		}

		name, err := m.GetPropertyName()
		if err != nil {
			return nil, err
		}

		f := template.ApiFieldTemplateData{
			FieldName:    util.SnakeCase(m.Name),
			FieldType:    ft,
			FieldTagName: name,
		}

		// fmt.Printf("m.Type: %s\n", m.Type.Name())

		if _, ok := m.Type.(spec.DefineStruct); ok {
			if msg, err := GenMessage(m.Type); err != nil {
				return nil, err
			} else {
				f.FieldMessage = msg
			}
		}

		if tag == "header" {
			if formatTag, err := GenFormatTag(m.Type); err != nil {
				return nil, err
			} else {
				f.FieldFormatTag = formatTag
			}
		}
		result = append(result, &f)
	}
	return result, nil
}

func GenMessage(t spec.Type) (*template.ApiMessageTemplateData, error) {
	result := template.ApiMessageTemplateData{
		MessageName: util.SnakeCase(t.Name()),
		Fields:      map[string][]*template.ApiFieldTemplateData{},
		FieldCount:  0,
	}

	definedType, ok := t.(spec.DefineStruct)
	if !ok {
		return nil, errors.New("no message of type " + t.Name())
	}

	if cJson, err := genCJson(definedType); err != nil {
		return nil, err
	} else {
		result.CJson = cJson
	}

	for _, tag := range requestFieldTags {
		if fields, err := GenFieldsByTag(definedType, tag); err != nil {
			return nil, err
		} else {
			result.Fields[tag] = fields
			result.FieldCount += len(fields)
		}
	}
	result.HeaderCount = len(result.Fields["header"])
	result.BodyCount = len(result.Fields["json"])
	result.PathCount = len(result.Fields["path"])
	result.FormCount = len(result.Fields["form"])

	return &result, nil
}

func GenMessages(dir string, api *spec.ApiSpec) error {
	data := template.ApiMessagesTemplateData{
		Messages: []*template.ApiMessageTemplateData{},
	}

	// 遍历消息类型
	for _, t := range api.Types {
		if m, err := GenMessage(t); err != nil {
			return err
		} else {
			data.Messages = append(data.Messages, m)
		}
	}

	if err := template.GenFile(dir, "message.h", template.ApiMessageHeader, data); err != nil {
		return err
	}

	if err := template.GenFile(dir, "message.c", template.ApiMessageSource, data); err != nil {
		return err
	}

	return nil
}
