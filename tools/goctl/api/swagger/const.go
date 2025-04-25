package swagger

const (
	tagHeader   = "header"
	tagPath     = "path"
	tagForm     = "form"
	tagJson     = "json"
	defFlag     = "default="
	enumFlag    = "options="
	rangeFlag   = "range="
	exampleFlag = "example="

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

	swaggerSecurityDefinitionBearerAuth = "BearerAuth"
	swaggerSecurityDefinitionName       = "Authorization"
	swaggerSecurityDefinitionIn         = "header"
)
