package gen

import (
	"errors"

	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/pygen/util"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)


func GenField(m spec.Member) (*template.PyFieldTemplateData, error) {
	keyName, err := m.GetPropertyName()
	if err != nil {
		return nil, err
	}
	tag := ""
	tags := m.Tags()
	if len(tags) > 0 {
		tag = m.Tags()[0].Key
	}

	field := &template.PyFieldTemplateData{
		FieldName:    util.SnakeCase(m.Name),
		FieldTag:     tag,
		FieldTagName: keyName,
	}
	return field, nil
}

func GenMessage(t spec.Type) (*template.PyMessageTemplateData, error) {
	message := &template.PyMessageTemplateData{
		MessageName: util.PascalCase(t.Name()),
		Fields:      []*template.PyFieldTemplateData{},
	}

	definedType, ok := t.(spec.DefineStruct)
	if !ok {
		return nil, errors.New("no message of type " + t.Name())
	}

	for _, m := range definedType.Members {
		f, err := GenField(m)
		if err != nil {
			return nil, err
		}
		message.Fields = append(message.Fields, f)
	}
	message.HeaderCount = len(definedType.GetTagMembers("header"))
	message.BodyCount = len(definedType.GetTagMembers("json"))
	message.PathCount = len(definedType.GetTagMembers("path"))
	message.FormCount = len(definedType.GetTagMembers("form"))

	return message, nil
}

func GenMessages(dir string, api *spec.ApiSpec) error {
	data := template.PyMessagesTemplateData{
		Messages: []*template.PyMessageTemplateData{},
	}

	// 遍历消息类型
	for _, t := range api.Types {
		if m, err := GenMessage(t); err != nil {
			return err
		} else {
			data.Messages = append(data.Messages, m)
		}
	}

	return template.GenFile(dir, "message.py", template.ApiMessage, data)
}
