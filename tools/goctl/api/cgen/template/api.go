package template

import (
	_ "embed"
)

//go:embed api_base_h.tpl
var ApiBaseHeader string

//go:embed api_base_c.tpl
var ApiBaseSource string

//go:embed api_client_h.tpl
var ApiClientHeader string

//go:embed api_client_c.tpl
var ApiClientSource string

//go:embed api_message_h.tpl
var ApiMessageHeader string

//go:embed api_message_c.tpl
var ApiMessageSource string

type ApiActionTemplateData struct {
	ActionName      string
	HttpMethod      string
	UrlPrefix       string
	UrlPath         string
	RequestMessage  *ApiMessageTemplateData
	ResponseMessage *ApiMessageTemplateData
}

type ApiClientTemplateData struct {
	ClientName string
	Actions    []*ApiActionTemplateData
}

type ApiFieldTemplateData struct {
	FieldName      string
	FieldType      string
	FieldTagName   string
	FieldFormatTag string
	FieldMessage   *ApiMessageTemplateData
}

type ApiMessageTemplateData struct {
	MessageName string
	Fields      map[string][]*ApiFieldTemplateData
	FieldCount  int
	HeaderCount int
	BodyCount   int
	FormCount   int
	PathCount   int
	CJson       *CJsonTemplateData
}

type ApiMessagesTemplateData struct {
	Messages []*ApiMessageTemplateData
}
