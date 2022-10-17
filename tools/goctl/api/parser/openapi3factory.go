package parser

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

type extensionPropsHandler func(schemaInfo *openApi3SpecTypeSchemaInfo) error

const (
	schemaRpcTypeDefault = iota
	schemaRpcTypeRequest
	schemaRpcTypeResponse
)

type specRpcType struct {
	spec.Type
	openApi3RpcTypeInfo
}

type openApi3SpecTypeFactory interface {
	initParentInfo(parentFactory openApi3SpecTypeFactory)
	getParentInfo() specTypeFactoryParentInfo

	initExtensionPropsFunc(initExtensionPropsFuncList []extensionPropsHandler)
	getExtensionPropsFunc() []extensionPropsHandler

	initSchemaInfo(propertyName string, schemaRef *openapi3.SchemaRef) error
	getSchemaInfo() openApi3SpecTypeSchemaInfo

	setSpecTypeName() error
	initSpecTypeName(name string)

	getSpecTypeName() string
	getSpecType() (spec.Type, []spec.Type, error)
}

type openApi3SpecTypeSchemaInfo struct {
	schema            *openapi3.SchemaRef
	ref               string
	propertyName      string
	extensionTypeName string
	nullable          bool
	comment           string
}

type specTypeFactoryParentInfo struct {
	parentFieldNames []string
	rootTypeName     string
}

type openApi3SpecTypeFactoryBase struct {
	openApi3SpecTypeSchemaInfo

	parentInfo                    specTypeFactoryParentInfo
	initExtensionPropsHandlerList []extensionPropsHandler
	typeName                      string
}

func (factory *openApi3SpecTypeFactoryBase) initExtensionPropsFunc(initExtensionPropsHandlerList []extensionPropsHandler) {
	factory.initExtensionPropsHandlerList = initExtensionPropsHandlerList
}

func (factory *openApi3SpecTypeFactoryBase) getExtensionPropsFunc() []extensionPropsHandler {
	return factory.initExtensionPropsHandlerList
}

func (factory *openApi3SpecTypeFactoryBase) initParentInfo(parent openApi3SpecTypeFactory) {
	if parent == nil {
		return
	}
	parentParentInfo := parent.getParentInfo()
	schemaInfo := parent.getSchemaInfo()
	factory.parentInfo.rootTypeName = parentParentInfo.rootTypeName
	if factory.parentInfo.rootTypeName == "" {
		factory.parentInfo.rootTypeName = parent.getSpecTypeName()
	}
	if schemaInfo.propertyName == "" {
		return
	}
	factory.parentInfo.parentFieldNames = append(parentParentInfo.parentFieldNames, schemaInfo.propertyName)
}

func (factory *openApi3SpecTypeFactoryBase) getParentInfo() specTypeFactoryParentInfo {
	return factory.parentInfo
}

func (factory *openApi3SpecTypeFactoryBase) initSchemaInfo(propertyName string, schemaRef *openapi3.SchemaRef) error {
	schemaInfo := openApi3SpecTypeSchemaInfo{}
	schemaInfo.schema = schemaRef
	schemaInfo.comment = strings.Replace(schemaRef.Value.Description, "\n", ", ", -1)
	schemaInfo.propertyName = propertyName
	schemaInfo.ref = schemaRef.Ref
	schemaInfo.nullable = schemaRef.Value.Nullable
	for _, initExtensionPropsFunc := range factory.getExtensionPropsFunc() {
		if err := initExtensionPropsFunc(&schemaInfo); err != nil {
			return err
		}
	}
	factory.openApi3SpecTypeSchemaInfo = schemaInfo
	return nil
}

func (factory *openApi3SpecTypeFactoryBase) getSchemaInfo() openApi3SpecTypeSchemaInfo {
	return factory.openApi3SpecTypeSchemaInfo
}

func (factory *openApi3SpecTypeFactoryBase) initSpecTypeName(name string) {
	factory.typeName = name
}

func (factory *openApi3SpecTypeFactoryBase) setSpecTypeName() error {
	if factory.typeName != "" {
		return nil
	}
	if factory.extensionTypeName != "" {
		factory.typeName = factory.extensionTypeName
		return nil
	}

	if splits := strings.Split(factory.ref, "/"); len(splits) > 0 {
		schemaName := splits[len(splits)-1]
		if factory.ref == fmt.Sprintf("%s%s", _componentSchemasName, schemaName) {
			factory.typeName = util.Title(schemaName)
			return nil
		}
	}
	return nil
}


func (factory *openApi3SpecTypeFactoryBase) getSpecTypeName() string {
	return factory.typeName
}

type openApi3RpcTypeInfo struct {
	operationId   string
	schemaRpcType int
	url           string
	method        string
	httpCode      string
	encodeType    string
}


type openApi3RpcTypeFactory struct {
	openApi3SpecTypeFactoryBase
	openApi3RpcTypeInfo

	subFactory openApi3SpecTypeFactory
}

func (factory *openApi3RpcTypeFactory) setSpecTypeName() error {
	if err := factory.openApi3SpecTypeFactoryBase.setSpecTypeName(); err != nil {
		return err
	}
	if name := factory.openApi3SpecTypeFactoryBase.getSpecTypeName(); name != "" {
		factory.typeName = util.Title(name)
		return nil
	}

	if factory.schemaRpcType == schemaRpcTypeRequest {
		if factory.operationId == "" {
			return fmt.Errorf("url: %s, method: %s, request type need name, you can define extension x-type-name or operateId ", factory.url, factory.method)
		}
		name := fmt.Sprintf("%sInput", factory.operationId)
		factory.typeName = util.Title(name)
		return nil
	}

	if factory.schemaRpcType == schemaRpcTypeResponse {
		if factory.operationId == "" {
			return fmt.Errorf("url: %s, method: %s, response type need name, you can define extension x-type-name or operateId", factory.url, factory.method)
		}
		name := fmt.Sprintf("%sOutput", factory.operationId)
		factory.typeName = util.Title(name)
		return nil
	}
	return nil
}

func (factory *openApi3RpcTypeFactory) getSpecType() (spec.Type, []spec.Type, error) {
	if err := factory.setSpecTypeName(); err != nil {
		return nil, nil, err
	}
	factory.subFactory.initSpecTypeName(factory.getSpecTypeName())

	specType, childSpecType, err := factory.subFactory.getSpecType()
	if err != nil {
		return nil, nil, err
	}
	specRpcType := specRpcType{
		Type:                specType,
		openApi3RpcTypeInfo: factory.openApi3RpcTypeInfo,
	}
	return specRpcType, childSpecType, nil
}

type openApi3SpecPrimitiveFactory struct {
	openApi3SpecTypeFactoryBase

	Min          *float64
	Max          *float64
	MultipleOf   *float64
	ExclusiveMin bool
	ExclusiveMax bool

	MinLength uint64
	MaxLength *uint64
	Pattern   string
	Format    string

	Enum []interface{}
}

func (factory *openApi3SpecPrimitiveFactory) setSpecTypeName() error {
	typeName, ok := primitiveTypeMap[factory.schema.Value.Type]
	if !ok {
		return fmt.Errorf("schema type: %s error", factory.schema.Value.Type)
	}
	factory.typeName = typeName
	return nil
}

func (factory *openApi3SpecPrimitiveFactory) getSpecType() (spec.Type, []spec.Type, error) {
	if err := factory.setSpecTypeName(); err != nil {
		return nil, nil, err
	}
	specType := spec.PrimitiveType{
		BaseType: spec.BaseType{
			RawNullable: factory.nullable,
		},
		Comment:      factory.comment,
		RawName:      factory.typeName,
		Min:          factory.schema.Value.Min,
		Max:          factory.schema.Value.Max,
		MultipleOf:   factory.schema.Value.MultipleOf,
		ExclusiveMin: factory.schema.Value.ExclusiveMin,
		ExclusiveMax: factory.schema.Value.ExclusiveMax,
		MinLength:    factory.schema.Value.MinLength,
		MaxLength:    factory.schema.Value.MaxLength,
		Pattern:      factory.schema.Value.Pattern,
		Format:       factory.schema.Value.Format,
		Enum:         factory.schema.Value.Enum,
	}
	return specType, nil, nil
}

type openApi3SpecObjectFactory struct {
	openApi3SpecTypeFactoryBase
	required []string
}

func (factory *openApi3SpecObjectFactory) setSpecTypeName() error {
	if err := factory.openApi3SpecTypeFactoryBase.setSpecTypeName(); err != nil {
		return err
	}
	if name := factory.openApi3SpecTypeFactoryBase.getSpecTypeName(); name != "" {
		factory.typeName = util.Title(name)
		return nil
	}

	if factory.parentInfo.rootTypeName == "" || factory.propertyName == "" {
		return errors.New("can not set type name")
	}
	fieldNames := append(factory.parentInfo.parentFieldNames, factory.propertyName)
	name := fmt.Sprintf("%s_%s", factory.parentInfo.rootTypeName, strings.Join(fieldNames, "_"))
	factory.typeName = util.Title(name)
	return nil
}

func (factory *openApi3SpecObjectFactory) getSpecType() (spec.Type, []spec.Type, error) {
	if err := factory.setSpecTypeName(); err != nil {
		return nil, nil, err
	}
	specType := spec.DefineStruct{
		BaseType: spec.BaseType{
			RawNullable: factory.nullable,
		},
		Comment:  factory.comment,
		Required: factory.schema.Value.Required,
		RawName:  factory.getSpecTypeName(),
	}
	childType := make([]spec.Type, 0)
	for schemaName, schemaRef := range factory.schema.Value.Properties {
		memberFactory, err := newOpenApi3SpecTypeFactory(schemaRef, factory.getExtensionPropsFunc(), parentFactoryOption(factory), propertyNameOption(schemaName))
		if err != nil {
			return nil, nil, err
		}
		memberType, memberChildType, err := memberFactory.getSpecType()
		if err != nil {
			return nil, nil, err
		}
		childType = append(childType, memberType)
		childType = append(childType, memberChildType...)
		isRequired := false
		for _, requireProperty := range factory.schema.Value.Required {
			if requireProperty == schemaName {
				isRequired = true
			}
		}
		member := spec.Member{
			Name:       schemaName,
			Type:       memberType,
			Comment:    strings.Join(memberType.Comments(), ","),
			IsRequired: isRequired,
		}

		jsonTag := getMemberTypeJsonTag(schemaName, isRequired)
		member.Tag = formatTags(jsonTag)
		specType.Members = append(specType.Members, member)
	}

	sort.Slice(specType.Members, func(i, j int) bool {
		return specType.Members[i].Name < specType.Members[j].Name
	})
	return specType, childType, nil
}

type openApi3SpecArrayFactory struct {
	openApi3SpecTypeFactoryBase
}

func (factory *openApi3SpecArrayFactory) setSpecTypeName() error {
	return nil
}

func (factory *openApi3SpecArrayFactory) getSpecType() (spec.Type, []spec.Type, error) {
	itemFactory, err := newOpenApi3SpecTypeFactory(factory.schema.Value.Items, factory.getExtensionPropsFunc(), parentFactoryOption(factory), propertyNameOption("item"))
	if err != nil {
		return nil, nil, err
	}

	itemType, itemChildType, err := itemFactory.getSpecType()
	if err != nil {
		return nil, nil, err
	}
	childType := append(itemChildType, itemType)
	factory.initSpecTypeName(fmt.Sprintf("[]%s", itemType.Name()))

	arrayType := spec.ArrayType{
		Value: itemType,
		BaseType: spec.BaseType{
			RawNullable: factory.nullable,
		},
		RawName:     factory.getSpecTypeName(),
		Comment:     factory.comment,
		MinItems:    factory.schema.Value.MinItems,
		MaxItems:    factory.schema.Value.MaxItems,
		UniqueItems: factory.schema.Value.UniqueItems,
	}

	return arrayType, childType, nil
}

type openApi3SpecMapFactory struct {
	openApi3SpecTypeFactoryBase
}

func (factory *openApi3SpecMapFactory) setSpecTypeName() error {
	return nil
}

func (factory *openApi3SpecMapFactory) getSpecType() (spec.Type, []spec.Type, error) {
	if factory.schema.Value == nil {
		return nil, nil, errors.New("schema type is null")
	}

	valueFactory, err := newOpenApi3SpecTypeFactory(factory.schema.Value.AdditionalProperties, factory.getExtensionPropsFunc(), parentFactoryOption(factory), propertyNameOption("value"))
	if err != nil {
		return nil, nil, err
	}

	valueType, valueChildType, err := valueFactory.getSpecType()
	if err != nil {
		return nil, nil, err
	}
	childType := append(valueChildType, valueType)
	factory.initSpecTypeName(fmt.Sprintf("map[string] %s", valueType.Name()))
	mapType := spec.MapType{
		Value: valueType,
		BaseType: spec.BaseType{
			RawNullable: factory.nullable,
		},
		RawName: factory.getSpecTypeName(),
		Comment: factory.comment,
		//openapi map的key只支持string
		Key: "string",
	}
	return mapType, childType, nil
}

var primitiveTypeMap = map[string]string{
	openapi3.TypeBoolean: "bool",
	openapi3.TypeNumber:  "float64",
	openapi3.TypeString:  "string",
	openapi3.TypeInteger: "int64",
	//type 没有定义openapi 中为string
	"": "string",
}

type openApi3SpecTypeFactoryConfig struct {
	propertyName        string
	openApi3RpcTypeInfo *openApi3RpcTypeInfo

	parentFactory openApi3SpecTypeFactory
}

type factoryConfigOption func(cfg *openApi3SpecTypeFactoryConfig)

func propertyNameOption(propertyName string) factoryConfigOption {
	return func(cfg *openApi3SpecTypeFactoryConfig) {
		cfg.propertyName = propertyName
	}
}

func rpcInfoOption(openApi3RpcTypeInfo *openApi3RpcTypeInfo) factoryConfigOption {
	return func(cfg *openApi3SpecTypeFactoryConfig) {
		cfg.openApi3RpcTypeInfo = openApi3RpcTypeInfo
	}
}

func parentFactoryOption(parentFactory openApi3SpecTypeFactory) factoryConfigOption {
	return func(cfg *openApi3SpecTypeFactoryConfig) {
		cfg.parentFactory = parentFactory
	}
}

func newOpenApi3SpecTypeFactory(schema *openapi3.SchemaRef, extensionPropsFuncList []extensionPropsHandler, options ...factoryConfigOption) (openApi3SpecTypeFactory, error) {
	if schema.Value == nil {
		return nil, errors.New("schema type is null")
	}
	cfg := openApi3SpecTypeFactoryConfig{}
	for _, option := range options {
		option(&cfg)
	}
	var factory openApi3SpecTypeFactory
	if cfg.openApi3RpcTypeInfo != nil {
		rpcFactory := &openApi3RpcTypeFactory{}
		options = append(options, rpcInfoOption(nil))
		subFactory, err := newOpenApi3SpecTypeFactory(schema, extensionPropsFuncList, options...)
		if err != nil {
			return nil, err
		}
		rpcFactory.openApi3RpcTypeInfo = *cfg.openApi3RpcTypeInfo
		rpcFactory.subFactory = subFactory
		factory = rpcFactory
	} else if schema.Value.Type == openapi3.TypeObject && schema.Value.AdditionalProperties == nil {
		//schema.Value.AdditionalProperties 有值表明是map结构
		factory = &openApi3SpecObjectFactory{}
	} else if schema.Value.Type == openapi3.TypeObject && schema.Value.AdditionalProperties != nil {
		factory = &openApi3SpecMapFactory{}
	} else if schema.Value.Type == openapi3.TypeArray {
		factory = &openApi3SpecArrayFactory{}
	} else if _, ok := primitiveTypeMap[schema.Value.Type]; ok {
		factory = &openApi3SpecPrimitiveFactory{}
	} else {
		return nil, fmt.Errorf("unexpect openapi3 type: %s", schema.Value.Type)
	}
	factory.initExtensionPropsFunc(extensionPropsFuncList)
	if err := factory.initSchemaInfo(cfg.propertyName, schema); err != nil {
		return nil, err
	}
	factory.initParentInfo(cfg.parentFactory)
	return factory, nil
}