package swagger

const (
	tagHeader    = "header"
	tagPath      = "path"
	tagForm      = "form"
	tagJson      = "json"
	defFlag      = "default="
	enumFlag     = "options="
	rangeFlag    = "range="
	exampleFlag  = "example="
	optionalFlag = "optional"

	paramsInHeader = "header"
	paramsInPath   = "path"
	paramsInQuery  = "query"
	paramsInBody   = "body"
	paramsInForm   = "formData"

	swaggerTypeInteger = "integer"
	swaggerTypeNumber  = "number"
	swaggerTypeString  = "string"
	swaggerTypeBoolean = "boolean"
	swaggerTypeArray   = "array"
	swaggerTypeObject  = "object"

	swaggerVersion  = "2.0"
	applicationJson = "application/json"
	applicationForm = "application/x-www-form-urlencoded"
	schemeHttps     = "https"
	defaultHost     = "127.0.0.1"
	defaultBasePath = "/"
)

const (
	propertyKeyUseDefinitions          = "useDefinitions"
	propertyKeyExternalDocsDescription = "externalDocsDescription"
	propertyKeyExternalDocsURL         = "externalDocsURL"
	propertyKeyDescription             = "description"
	propertyKeyProduces                = "produces"
	propertyKeyConsumes                = "consumes"
	propertyKeySchemes                 = "schemes"
	propertyKeyTags                    = "tags"
	propertyKeySummary                 = "summary"
	propertyKeyDeprecated              = "deprecated"
	propertyKeyPrefix                  = "prefix"
	propertyKeyAuthType                = "authType"
	propertyKeyHost                    = "host"
	propertyKeyBasePath                = "basePath"
	propertyKeyWrapCodeMsg             = "wrapCodeMsg"
	propertyKeyBizCodeEnumDescription  = "bizCodeEnumDescription"
)

const (
	defaultValueOfPropertyUseDefinition = false
)
