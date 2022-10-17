package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

const (
	_componentSchemasName = "#/components/schemas/"
	_jsonContent          = "application/json"
	_serverNameKey        = "server-name"
	_httpCodeDefault      = "default"
	_groupKey             = "group"
	openApi               = "openApi"
)

type openApi3RpcTypes struct {
	specTypes []spec.Type
}

type openApi3Parser struct {
	fileName               string
	doc                    *openapi3.T
	spec                   *spec.ApiSpec
	extensionPropsFuncList []extensionPropsHandler
	openApi3RpcTypes       *openApi3RpcTypes
}

func newOpenApi3Parser(ctx context.Context, filename string) (openApi3Parser, error) {
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	openapi3Doc, err := loader.LoadFromFile(filename)
	if err != nil {
		return openApi3Parser{}, err
	}

	return newOpenApi3ParserV2(ctx, filename, openapi3Doc)
}

func newOpenApi3ParserV2(ctx context.Context, filename string, openapi3Doc *openapi3.T) (openApi3Parser, error) {
	err := openapi3Doc.Validate(ctx)
	if err != nil {
		return openApi3Parser{}, err
	}
	openApi3Parser := openApi3Parser{doc: openapi3Doc, spec: &spec.ApiSpec{}, fileName: filename, openApi3RpcTypes: &openApi3RpcTypes{}}
	openApi3Parser.extensionPropsFuncList = append(openApi3Parser.extensionPropsFuncList, setSchemaExtensionPropsTypeName)
	return openApi3Parser, nil
}

func setSchemaExtensionPropsTypeName(schemaInfo *openApi3SpecTypeSchemaInfo) error {
	if schemaInfo.schema == nil || schemaInfo.schema.Value == nil {
		return nil
	}
	typeName, ok := schemaInfo.schema.Value.Extensions["x-typename"]
	if !ok {
		return nil
	}
	extensionTypeNameJsonBytes, ok := typeName.(json.RawMessage)
	if !ok {
		return fmt.Errorf("x-typename err type, expect: json.RawMessage, actual: %T", typeName)
	}

	err := json.Unmarshal(extensionTypeNameJsonBytes, &schemaInfo.extensionTypeName)
	if err != nil {
		return fmt.Errorf("x-nullable value failed, should been true or false")
	}
	return nil
}

func (p openApi3Parser) parse() (*spec.ApiSpec, error) {
	if p.doc == nil {
		return nil, errors.New("openapi3 doc is null")
	}

	if err := p.validateOpenApiDoc(); err != nil {
		return nil, err
	}
	if err := p.convert2Spec(); err != nil {
		return nil, err
	}

	return p.spec, nil
}

func (p openApi3Parser) validateOpenApiDoc() error {
	for url, pathInfo := range p.doc.Paths {
		for method, operation := range pathInfo.Operations() {
			if operation == nil {
				continue
			}
			if operation.OperationID == "" {
				return fmt.Errorf("url: %s, method: %s, operateId is null, hander Method name depend operateId", url, method)
			}
		}
	}
	return nil
}

func (p openApi3Parser) convert2Spec() error {
	p.fillInfo()
	p.fillSyntax()
	err := p.findTypes()
	if err != nil {
		return err
	}
	err = p.fillService()
	if err != nil {
		return err
	}
	p.filterToGenTypes()

	return nil
}

func (p openApi3Parser) fillInfo() {
	if p.doc.Info == nil {
		return
	}
	properties := make(map[string]string)
	for key, value := range p.doc.Info.Extensions {
		properties[key] = fmt.Sprint(value)
	}
	properties["Title"] = p.doc.Info.Title
	properties["Desc"] = p.doc.Info.Description
	properties["Version"] = p.doc.Info.Version
	if p.doc.Info.Contact != nil {
		properties["Author"] = p.doc.Info.Contact.Name
		properties["Email"] = p.doc.Info.Contact.Email
	}

	p.spec.Info.Properties = properties
}

func (p openApi3Parser) fillSyntax() {
	p.spec.Syntax.Version = fmt.Sprintf("%s-%s", openApi, p.doc.OpenAPI)
	p.spec.Syntax.Comment = []string{p.doc.Info.Description}
}

// getResponseDefineName 更加openapi operationID获取返回的结构体名称
func (p openApi3Parser) getResponseTypeName(operationID string) string {
	return fmt.Sprintf("%sOutput", operationID)
}

func (p openApi3Parser) getRequestTypeName(operationID string) string {
	return fmt.Sprintf("%sInput", operationID)
}

func (p openApi3Parser) filterToGenTypes() {
	tmpTypes := make([]spec.Type, 0, len(p.spec.Types))
	for _, specType := range p.spec.Types {
		if structType, ok := specType.(spec.DefineStruct); ok {
			tmpTypes = append(tmpTypes, structType)
		}
	}
	p.spec.Types = tmpTypes
}

func (p openApi3Parser) findTypes() error {
	specTypes := make([]spec.Type, 0)
	requestTypes, requestChildTypes, err := p.findTypesRequestType()
	if err != nil {
		return err
	}
	p.openApi3RpcTypes.specTypes = append(p.openApi3RpcTypes.specTypes, requestTypes...)
	specTypes = append(specTypes, requestChildTypes...)
	responseTypes, responseChildTypes, err := p.findTypesResponseType()
	if err != nil {
		return err
	}
	p.openApi3RpcTypes.specTypes = append(p.openApi3RpcTypes.specTypes, responseTypes...)
	specTypes = append(specTypes, responseChildTypes...)
	for _, specType := range p.openApi3RpcTypes.specTypes {
		if rpcType, ok := specType.(specRpcType); ok {
			specTypes = append(specTypes, rpcType.Type)
		}
	}

	filterRepeatFunc := func(specTypes []spec.Type, checkedType spec.Type) ([]spec.Type, error) {
		for _, specType := range specTypes {
			if reflect.DeepEqual(specType, checkedType) {
				return specTypes, nil
			}
			if specType.Name() == checkedType.Name() {
				return nil, fmt.Errorf("type name: %s is repeate", specType.Name())
			}
		}
		specTypes = append(specTypes, checkedType)
		return specTypes, nil
	}
	retSpecTypes := make([]spec.Type, 0)
	for _, specType := range specTypes {
		if _, ok := specType.(spec.DefineStruct); !ok {
			continue
		}
		retSpecTypes, err = filterRepeatFunc(retSpecTypes, specType)
		if err != nil {
			return err
		}
	}
	p.spec.Types = retSpecTypes

	return nil
}

func (p openApi3Parser) findTypesRequestType() ([]spec.Type, []spec.Type, error) {
	childTypes := make([]spec.Type, 0)
	requestTypes := make([]spec.Type, 0)
	for url, pathInfo := range p.doc.Paths {
		for method, operation := range pathInfo.Operations() {
			if operation == nil {
				continue
			}
			paramSpecType, paramChildTypes, err := p.findTypeBySchemaParams(operation.OperationID, url, method, operation.Parameters)
			if err != nil {
				return nil, nil, fmt.Errorf("url: %s, method: %s get request param type failed: %w", url, method, err)
			}
			childTypes = append(childTypes, paramChildTypes...)

			requestSpecTypeList, requestBodyChildTypes, err := p.findTypeBySchemaRequestBody(operation.OperationID, url, method, operation.RequestBody)
			if err != nil {
				return nil, nil, fmt.Errorf("url: %s, method: %s get requestBody type failed: %w", url, method, err)
			}
			childTypes = append(childTypes, requestBodyChildTypes...)
			if len(requestSpecTypeList) == 0 {
				requestTypes = append(requestTypes, paramSpecType)
				continue
			}
			if paramSpecType == nil {
				requestTypes = append(requestTypes, requestSpecTypeList...)
				continue
			}
			for _, requestSpecType := range requestSpecTypeList {
				specType, err := p.mergeRequestBodyAndParam(requestSpecType, paramSpecType)
				if err != nil {
					return nil, nil, fmt.Errorf("url: %s, method: %s meger request body and param failed: %w", url, method, err)
				}
				requestTypes = append(requestTypes, specType)
			}
		}
	}

	return requestTypes, childTypes, nil
}

func (p openApi3Parser) mergeRequestBodyAndParam(requestBody spec.Type, param spec.Type) (spec.Type, error) {
	requestSpecRpcType, ok := requestBody.(specRpcType)
	if !ok {
		return nil, fmt.Errorf("requestBody type: %s expect: specRpcType, actual: %T", requestBody.Name(), requestBody)
	}
	requestDefineStruct, ok := requestSpecRpcType.Type.(spec.DefineStruct)
	if !ok {
		return nil, fmt.Errorf("requestBody type: %s expect: specRpcType, actual: %T", requestBody.Name(), requestBody)
	}
	paramSpecRpcType, ok := param.(specRpcType)
	if !ok {
		return nil, fmt.Errorf("request param type: %s expect: specRpcType, actual: %T", requestBody.Name(), requestBody)
	}
	paramDefineStruct, ok := paramSpecRpcType.Type.(spec.DefineStruct)
	if !ok {
		return nil, fmt.Errorf("request param type: %s expect: specRpcType, actual: %T", requestBody.Name(), requestBody)
	}
	for _, requestBodyMember := range requestDefineStruct.Members {
		for _, paramMember := range paramDefineStruct.Members {
			if requestBodyMember.Name == paramMember.Name {
				return nil, fmt.Errorf("field name: %s borth in requestBody and param", requestBodyMember.Name)
			}
		}
	}
	requestDefineStruct.Members = append(requestDefineStruct.Members, paramDefineStruct.Members...)
	requestDefineStruct.Required = append(requestDefineStruct.Required, paramDefineStruct.Required...)
	requestSpecRpcType.Type = requestDefineStruct
	return requestSpecRpcType, nil
}

func (p openApi3Parser) findTypesResponseType() ([]spec.Type, []spec.Type, error) {
	specTypes := make([]spec.Type, 0)
	childTypes := make([]spec.Type, 0)
	for url, pathInfo := range p.doc.Paths {
		for method, operation := range pathInfo.Operations() {
			if operation == nil {
				continue
			}
			responseSpecTypeList, responseChildTypes, err := p.findTypeBySchemaResponses(operation.OperationID, url, method, operation.Responses)
			if err != nil {
				return nil, nil, err
			}
			specTypes = append(specTypes, responseSpecTypeList...)
			childTypes = append(childTypes, responseChildTypes...)
		}
	}
	return specTypes, childTypes, nil
}

func (p openApi3Parser) findTypeBySchemaResponses(operationID, url, method string, responses openapi3.Responses) ([]spec.Type, []spec.Type, error) {
	if responses == nil {
		return nil, nil, nil
	}
	resSpecTypeList := make([]spec.Type, 0)
	childTypeList := make([]spec.Type, 0)
	for httpCode, responseInfo := range responses {
		if responseInfo.Value == nil {
			return nil, nil, fmt.Errorf("operationID: %s, no response value", operationID)
		}
		if responseInfo.Value.Content == nil {
			return nil, nil, fmt.Errorf("operationID: %s not exist, repsponse body", operationID)
		}
		for encodeType, responseSchemaInfo := range responseInfo.Value.Content {
			if responseSchemaInfo.Schema == nil {
				return nil, nil, fmt.Errorf("operationID: %s, not exist json encode Schema reponse define", operationID)
			}
			openApi3RpcTypeInfo := openApi3RpcTypeInfo{
				operationId:   operationID,
				schemaRpcType: schemaRpcTypeResponse,
				url:           url,
				method:        method,
				httpCode:      httpCode,
				encodeType:    encodeType,
			}
			factory, err := newOpenApi3SpecTypeFactory(responseSchemaInfo.Schema, p.extensionPropsFuncList, rpcInfoOption(&openApi3RpcTypeInfo))
			specType, childSpecType, err := factory.getSpecType()
			if err != nil {
				return nil, nil, err
			}
			resSpecTypeList = append(resSpecTypeList, specType)
			childTypeList = append(childTypeList, childSpecType...)
		}
	}
	return resSpecTypeList, childTypeList, nil
}

func (p openApi3Parser) findTypeBySchemaRequestBody(operationID, url, method string, requestBodyInfo *openapi3.RequestBodyRef) ([]spec.Type, []spec.Type, error) {
	if requestBodyInfo == nil {
		return nil, nil, nil
	}
	childTypes := make([]spec.Type, 0)
	resSpecTypeList := make([]spec.Type, 0)
	for encodeType, responseSchemaInfo := range requestBodyInfo.Value.Content {
		if responseSchemaInfo.Schema == nil {
			return nil, nil, fmt.Errorf("operationID: %s, not exist json encode Schema request define", operationID)
		}
		openApi3RpcTypeInfo := openApi3RpcTypeInfo{
			operationId:   operationID,
			schemaRpcType: schemaRpcTypeRequest,
			url:           url,
			method:        method,
			encodeType:    encodeType,
		}
		factory, err := newOpenApi3SpecTypeFactory(responseSchemaInfo.Schema, p.extensionPropsFuncList, rpcInfoOption(&openApi3RpcTypeInfo))
		specType, childSpecType, err := factory.getSpecType()
		if err != nil {
			return nil, nil, err
		}
		resSpecTypeList = append(resSpecTypeList, specType)
		childTypes = append(childTypes, childSpecType...)
	}
	return resSpecTypeList, childTypes, nil
}

func getMemberTypeTag(tagKey string, tagValue string) string {
	return fmt.Sprintf(`%s:"%s"`, tagKey, tagValue)
}

func getMemberTypeJsonTag(name string, required bool) string {
	if !required {
		name = fmt.Sprintf("%s,omitempty", name)
	}
	return getMemberTypeTag("json", name)
}

func formatTags(tags ...string) string {
	tagValue := strings.Join(tags, " ")
	return fmt.Sprintf("`%s`", tagValue)
}

// getMemberTypeJsonTag 获取spec结构体成员变量的form标签
func getMemberTypeFormTag(name string) string {
	return getMemberTypeTag("form", name)
}

func getMemberTypePathTag(name string) string {
	return getMemberTypeTag("path", name)
}


func (p openApi3Parser) getTagsBySchemaParams(paramRef *openapi3.ParameterRef) (string, error) {
	if paramRef.Ref != "" {
		return "", nil
	}
	if paramRef.Value == nil {
		return "", nil
	}
	if paramRef.Value.In == "query" {
		return getMemberTypeFormTag(paramRef.Value.Name), nil
	}
	if paramRef.Value.In == "path" {
		return getMemberTypePathTag(paramRef.Value.Name), nil
	}
	return "", fmt.Errorf("param: %s input type: %s err, should been in or query", paramRef.Value.Name, paramRef.Value.In)
}

// findTypeBySchemaRequestBody 从openapi3 的Parameters对象中获取请求数据的结构体类型
func (p openApi3Parser) findTypeBySchemaParams(operationID, url, method string, params openapi3.Parameters) (spec.Type, []spec.Type, error) {
	if len(params) == 0 {
		return nil, nil, nil
	}
	rawName := p.getRequestTypeName(operationID)
	schemaInfo := openApi3SpecTypeSchemaInfo{}
	setSchemaExtensionPropsTypeName(&schemaInfo)
	if schemaInfo.extensionTypeName != "" {
		rawName = schemaInfo.extensionTypeName
	}
	if rawName == "" {
		return nil, nil, fmt.Errorf("url: %s, method: %s, param must define x-typename or operationId as param object name", url, method)
	}
	specType := spec.DefineStruct{
		RawName: rawName,
	}
	childTypes := make([]spec.Type, 0)
	for _, param := range params {
		if param.Value == nil {
			return spec.DefineStruct{}, nil, fmt.Errorf("url: %s, method: %s param is null", url, method)
		}

		paramName := param.Value.Name
		comment := param.Value.Description

		factory, err := newOpenApi3SpecTypeFactory(param.Value.Schema, p.extensionPropsFuncList, propertyNameOption(paramName))
		paramSpecType, childSpecTypes, err := factory.getSpecType()
		if err != nil {
			return nil, nil, err
		}
		childTypes = append(childTypes, childSpecTypes...)
		paramMember := spec.Member{
			Name:               paramName,
			Type:               paramSpecType,
			Comment:            comment,
			NotAllowEmptyValue: !param.Value.AllowEmptyValue,
		}
		if param.Value != nil {
			paramMember.IsRequired = param.Value.Required
		}

		tag, err := p.getTagsBySchemaParams(param)
		if err != nil {
			return nil, nil, err
		}
		paramMember.Tag = formatTags(tag)
		specType.Members = append(specType.Members, paramMember)
	}
	openApi3RpcTypeInfo := openApi3RpcTypeInfo{
		operationId:   operationID,
		schemaRpcType: schemaRpcTypeRequest,
		url:           url,
		method:        method,
	}
	retType := specRpcType{
		Type:                specType,
		openApi3RpcTypeInfo: openApi3RpcTypeInfo,
	}
	return retType, childTypes, nil

}

func (p openApi3Parser) getRequestType(url, method string) spec.Type {
	funcGetReq := func(encodeType string) spec.Type {
		for _, specType := range p.openApi3RpcTypes.specTypes {
			specRpcType, ok := specType.(specRpcType)
			if !ok {
				continue
			}
			if specRpcType.schemaRpcType != schemaRpcTypeRequest {
				continue
			}
			if specRpcType.url == url && specRpcType.method == method && specRpcType.encodeType == encodeType {
				return specRpcType.Type
			}
		}
		return nil
	}

	specType := funcGetReq(_jsonContent)
	if specType == nil {
		specType = funcGetReq("")
	}
	return specType
}

//getResponseType 根据operationID从已经解析的所有spec.Type 获取其对应返回数据类型
func (p openApi3Parser) getResponseType(url, method string) spec.Type {
	funcGetReq := func(httpCode, encodeType string) spec.Type {
		for _, specType := range p.openApi3RpcTypes.specTypes {
			specRpcType, ok := specType.(specRpcType)
			if !ok {
				continue
			}
			if specRpcType.schemaRpcType != schemaRpcTypeResponse {
				continue
			}
			if specRpcType.url == url && specRpcType.method == method &&
				specRpcType.encodeType == encodeType && specRpcType.httpCode == httpCode {
				return specRpcType.Type
			}
		}
		return nil
	}

	specType := funcGetReq(fmt.Sprint(http.StatusOK), _jsonContent)
	if specType != nil {
		return specType
	}

	specType = funcGetReq(_httpCodeDefault, _jsonContent)
	if specType != nil {
		return specType
	}
	return nil
}

func (p openApi3Parser) fixedPathUrl(url string, specType spec.Type) (string, error) {
	if specType == nil {
		return url, nil
	}
	structType, ok := specType.(spec.DefineStruct)
	if !ok {
		return url, nil
	}
	for _, member := range structType.Members {
		if !strings.Contains(member.Tag, getMemberTypePathTag(member.Name)) {
			continue
		}
		urlKey := fmt.Sprintf("{%s}", member.Name)
		if !strings.Contains(url, urlKey) {
			return "", fmt.Errorf("url: %s not contain field: %s", url, urlKey)
		}
		newUrlKey := fmt.Sprintf(":%s", member.Name)
		url = strings.Replace(url, urlKey, newUrlKey, 1)
	}
	return url, nil
}

func (p openApi3Parser) fillService() error {
	routeGroup, err := p.getRoute(p.doc.Paths)
	if err != nil {
		return err
	}
	annotation := spec.Annotation{
		Properties: make(map[string]string),
	}
	if p.doc.ExtensionProps.Extensions != nil {
		for key, value := range p.doc.ExtensionProps.Extensions {
			annotation.Properties[key] = fmt.Sprint(value)
		}
	}
	routeGroup.Annotation = annotation
	p.spec.Service.Groups = append(p.spec.Service.Groups, routeGroup)

	serverNameTag := p.doc.Tags.Get(_serverNameKey)
	if serverNameTag != nil {
		p.spec.Service.Name = serverNameTag.Description
	} else {
		fileName := path.Base(p.fileName)
		fileSuffix := path.Ext(p.fileName)
		filePrefix := fileName[0 : len(fileName)-len(fileSuffix)]
		p.spec.Service.Name = filePrefix
	}
	return nil
}


func (p openApi3Parser) getRoute(paths openapi3.Paths) (spec.Group, error) {
	group := spec.Group{}
	for url, pathInfo := range paths {
		for method, operation := range pathInfo.Operations() {
			if operation == nil {
				continue
			}
			annotation := spec.Annotation{
				Properties: make(map[string]string),
			}

			for key, value := range operation.ExtensionProps.Extensions {
				annotation.Properties[key] = fmt.Sprint(value)
			}
			if len(operation.Tags) != 0 {
				annotation.Properties[_groupKey] = operation.Tags[0]
			}
			route := spec.Route{
				AtServerAnnotation: annotation,
				Method:             strings.ToLower(method),
				Handler:            operation.OperationID,
				HandlerDoc:         []string{operation.Description},
			}
			route.RequestType = p.getRequestType(url, method)
			route.ResponseType = p.getResponseType(url, method)
			newUrl, err := p.fixedPathUrl(url, route.RequestType)
			if err != nil {
				return group, err
			}
			route.Path = newUrl
			group.Routes = append(group.Routes, route)
		}
	}

	return group, nil
}
