package template

import _ "embed"

//go:embed to_json_primitive.tpl
var cJsonToPrimitiveTemplate string

//go:embed to_json_array.tpl
var cJsonToArrayTemplate string

//go:embed to_json_object.tpl
var cJsonToObjectTemplate string

//go:embed from_json_primitive.tpl
var cJsonFromPrimitiveTemplate string

//go:embed from_json_array.tpl
var cJsonFromArrayTemplate string

//go:embed from_json_object.tpl
var cJsonFromObjectTemplate string

type CJsonType string
type CJsonCtor string
type CJsonCheck string
type CJsonValue string

const (
	CJsonUnsupported CJsonType = "unsupported"
	CJsonPrimitive   CJsonType = "primitive"
	CJsonArray       CJsonType = "array"
	CJsonObject      CJsonType = "object"

	CJsonCreateString CJsonCtor = "cJSON_CreateString"
	CJsonCreateNumber CJsonCtor = "cJSON_CreateNumber"
	CJsonCreateBool   CJsonCtor = "cJSON_CreateBool"

	CJsonIsBool   CJsonCheck = "cJSON_IsBool"
	CJsonIsNumber CJsonCheck = "cJSON_IsNumber"
	CJsonIsString CJsonCheck = "cJSON_IsString"

	CJsonValueInt    CJsonValue = "valueint"
	CJsonValueDouble CJsonValue = "valuedouble"
	CJsonValueString CJsonValue = "valuestring"
)

type CJsonTemplateData struct {
	Indent   int
	VarType  CJsonType
	VarName  string
	VarExpr  string
	VarCType string
	VarCtor  CJsonCtor
	VarCheck CJsonCheck
	VarValue CJsonValue
	Items    *CJsonTemplateData
	Pairs    map[string]*CJsonTemplateData
}
