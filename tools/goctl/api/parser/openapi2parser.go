package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/invopop/yaml"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type openApi2Parser struct {
	fileName string
	doc      *openapi2.T
	spec     *spec.ApiSpec
}

func newOpenApi2Parser(filename string) (openApi2Parser, error) {
	inputYAML, err := ioutil.ReadFile(filename)
	if err != nil {
		return openApi2Parser{}, fmt.Errorf("read openapi2.0 file failed: %w", err)
	}
	var openapi2Doc openapi2.T
	if err = yaml.Unmarshal(inputYAML, &openapi2Doc); err != nil {
		return openApi2Parser{}, fmt.Errorf("unmarshal openapi2.0 file failed: %w", err)
	}
	return openApi2Parser{doc: &openapi2Doc, spec: &spec.ApiSpec{}, fileName: filename}, nil
}

func (p openApi2Parser) parse() (*spec.ApiSpec, error) {
	if p.doc == nil {
		return nil, errors.New("openapi2 doc is null")
	}

	openapi3Doc, err := openapi2conv.ToV3(p.doc)
	if err != nil {
		return nil, fmt.Errorf("conv to openapi 3.0 failed: %w", err)
	}
	openApi3Parser, err := newOpenApi3ParserV2(context.Background(), p.fileName, openapi3Doc)
	if err != nil {
		return nil, fmt.Errorf("conv to openapi 3.0 failed: %w", err)
	}
	openApi3Parser.extensionPropsFuncList = append(openApi3Parser.extensionPropsFuncList, dealNullableExtensionProps)
	return openApi3Parser.parse()
}

//dealBySchemaExtensionProps 从openapi3 的SchemaRef对象中获取结构体类型
func dealNullableExtensionProps(schemaInfo *openApi3SpecTypeSchemaInfo) error {
	if schemaInfo.schema.Value == nil {
		return fmt.Errorf("dealNullableExtensionProps failed: schema is null")
	}
	nullableInterface, ok := schemaInfo.schema.Value.Extensions["x-nullable"]
	if !ok {
		return nil
	}
	nullableJson, ok := nullableInterface.(json.RawMessage)
	if !ok {
		return fmt.Errorf("x-nullable err type, expect: json.RawMessage, actual: %T", nullableInterface)
	}
	var nullable bool
	err := json.Unmarshal(nullableJson, &nullable)
	if err != nil {
		return fmt.Errorf("x-nullable value failed, should been true or false")
	}
	schemaInfo.nullable = nullable
	return nil
}
