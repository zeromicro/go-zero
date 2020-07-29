package gen

import (
	"strings"

	"zero/tools/goctl/model/mongomodel/utils"
)

func genMethodTemplate(funcDesc FunctionDesc, needCache bool) (template string) {
	var tmp string
	switch funcDesc.Type {
	case functionTypeGet:
		if needCache {
			tmp = getTemplate
		} else {
			tmp = noCacheGetTemplate
		}
	case functionTypeFind:
		tmp = findTemplate
	case functionTypeSet:
		if needCache {
			tmp = ""
		} else {
			tmp = noCacheSetFieldtemplate
		}
	default:
		return ""
	}
	tmp = strings.ReplaceAll(tmp, "{{.Name}}", funcDesc.FieldName)
	tmp = strings.ReplaceAll(tmp, "{{.name}}", utils.UpperCamelToLower(funcDesc.FieldName))
	tmp = strings.ReplaceAll(tmp, "{{.type}}", funcDesc.FieldType)
	return tmp
}
