package gen

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/template"
	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/util"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func cJsonVarCtor(t spec.Type) template.CJsonCtor {
	switch t.Name() {
	case "bool":
		return template.CJsonCreateBool
	case "int", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float", "float32", "float64":
		return template.CJsonCreateNumber
	case "string":
		return template.CJsonCreateString
	default:
		return ""
	}
}

func cJsonVarCheckAndValue(t spec.Type) (template.CJsonCheck, template.CJsonValue) {
	switch t.Name() {
	case "bool":
		return template.CJsonIsBool, template.CJsonValueInt
	case "int", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
		return template.CJsonIsNumber, template.CJsonValueInt
	case "float", "float32", "float64":
		return template.CJsonIsNumber, template.CJsonValueDouble
	case "string":
		return template.CJsonIsString, template.CJsonValueString
	default:
		return "", ""
	}
}

func genCJsonPrimitive(t spec.PrimitiveType, varName string, varExpr string, indent int) (*template.CJsonTemplateData, error) {
	varCtor := cJsonVarCtor(t)
	if varCtor == "" {
		return nil, fmt.Errorf("no supported of cjson primitive type %s ", t.Name())
	}
	varCheck, varValue := cJsonVarCheckAndValue(t)
	if varCheck == "" {
		return nil, fmt.Errorf("no supported of cjson primitive type %s ", t.Name())
	}

	varCType, err := primitiveType(t.Name())
	if err != nil {
		return nil, err
	}

	primitive := &template.CJsonTemplateData{
		Indent:   indent,
		VarType:  template.CJsonPrimitive,
		VarName:  varName,
		VarExpr:  varExpr,
		VarCType: varCType,
		VarCtor:  varCtor,
		VarCheck: varCheck,
		VarValue: varValue,
	}

	return primitive, nil
}

func genCJsonItemsType(t spec.Type) (template.CJsonType, string, error) {
	switch vt := t.(type) {
	case spec.PrimitiveType:
		n, err := primitiveType(t.Name())
		return template.CJsonPrimitive, n, err
	case spec.PointerType:
		if t, n, err := genCJsonItemsType(vt.Type); err != nil {
			return template.CJsonUnsupported, "", err
		} else {
			return t, fmt.Sprintf("%s*", n), nil
		}
	case spec.DefineStruct:
		return template.CJsonObject, structType(t.Name()), nil
	}

	return template.CJsonUnsupported, "", fmt.Errorf("no supported of cjson items type %s ", t.Name())
}

func genCJsonArray(t spec.ArrayType, varName string, varExpr string, indent int) (*template.CJsonTemplateData, error) {
	itemType, itemCType, err := genCJsonItemsType(t.Value)
	if err != nil {
		return nil, err
	}
	itemName := fmt.Sprintf("%s_item", varName)
	itemExpr := fmt.Sprintf("v_%s_items[i]", varName)
	itemCtor := cJsonVarCtor(t.Value)
	itemsTemplateData := &template.CJsonTemplateData{
		Indent:   indent + 4,
		VarType:  itemType,
		VarName:  itemName,
		VarCtor:  itemCtor,
		VarCType: itemCType,
		VarExpr:  itemExpr,
	}

	if itemCtor == "" {
		ot := findTypeByName(t.Value.Name())
		if ot != nil {
			if vt, ok := ot.(spec.DefineStruct); ok {
				itemsTemplateData, err = genCJsonObject(vt, itemName, itemExpr, indent+4)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	array := &template.CJsonTemplateData{
		Indent:  indent,
		VarType: template.CJsonArray,
		VarName: varName,
		VarExpr: varExpr,
		// VarCType: ,
		Items: itemsTemplateData,
	}

	return array, nil
}

func genCJsonPair(t spec.Type, varName string, varExpr string, indent int) (*template.CJsonTemplateData, error) {
	switch vt := t.(type) {
	case spec.PrimitiveType:
		return genCJsonPrimitive(vt, varName, varExpr, indent)
	case spec.ArrayType:
		return genCJsonArray(vt, varName, varExpr, indent)
	case spec.DefineStruct:
		return genCJsonObject(vt, varName, varExpr, indent)
	}

	return nil, fmt.Errorf("no supported of cjson type: %s ", t.Name())
}

func genCJsonObject(definedType spec.DefineStruct, varName string, varExpr string, indent int) (*template.CJsonTemplateData, error) {
	object := &template.CJsonTemplateData{
		Indent:   indent,
		VarType:  template.CJsonObject,
		VarName:  varName,
		VarExpr:  varExpr,
		VarCType: fmt.Sprintf("%s_t", util.SnakeCase(definedType.Name())),
		Pairs:    map[string]*template.CJsonTemplateData{},
	}

	for _, m := range definedType.GetTagMembers("json") {
		pn, err := m.GetPropertyName()
		if err != nil {
			return nil, err
		}
		n := util.SnakeCase(m.Name)
		pair, err := genCJsonPair(m.Type, n, fmt.Sprintf("%s.%s", varExpr, n), indent+4)
		if err != nil {
			return nil, err
		}
		object.Pairs[pn] = pair
	}

	return object, nil
}

func genCJson(definedType spec.DefineStruct) (*template.CJsonTemplateData, error) {
	return genCJsonObject(definedType, "root", "(*message)", 4)
}
