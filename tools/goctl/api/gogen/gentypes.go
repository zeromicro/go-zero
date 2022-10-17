package gogen

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const typesFile = "types"

var (
	//go:embed types.tpl
	typesTemplate string

	//go:embed typefieldvalidator.tpl
	typeFieldValidateTemplate string

	//go:embed typevalidator.tpl
	typeFuncValidateTemplate string
)

var (
	//生成校验方法模板
	typeFuncValidatorTpl = template.Must(template.New("typeFuncValidate").Parse(typeFuncValidateTemplate))
	//生成成员变量校验模板
	typeFieldValidatorTpl = template.Must(template.New("typeFieldValidate").Parse(typeFieldValidateTemplate))
)

var (
	_validateFuncTypeAlias = "v"
	_packageFmt            = `"fmt"`
	_packageRegexp         = `"regexp"`
)

//生成type文件的生成器
type typeBuilder struct {
	//package 内容
	importPackagesContents []string
	//全员变量
	globalVarsContents []string
	//结构体类型
	typesContent string
	//结构体类型校验方法
	typeFuncValidatorsContents []string
	//需要生成结构体类型
	types []spec.Type
}

// typeValidateFuncBuilder 生成结构体类型检验方法的生成器
type typeValidateFuncBuilder struct {
	typeGen               *typeBuilder // 生成type文件的生成器
	localVarsContents     []string     // 方法中的变量内容
	propertyNullContents  []string     // 需要进行空校验的内容
	fieldValidateContents []string     // 结构体成员校验内容
}

// typeFieldValidateGen 生成结构体成员属性校验的生成器
type typeFieldValidateBuilder struct {
	*typeValidateFuncBuilder          // 生成结构体类型检验方法的生成器
	typeAlias                string   // 结构体类型方法接受变量别名
	isRequire                bool     // 当前变量是否为必传
	rootTypeName             string   // 根类型名
	parentFieldNames         []string // 包括父节点及父父节点在其类型中成员名称，全局变量regex正则的名称由  rootTypeName...parentFieldNames.fieldName 组成
	fieldName                string   // 成员变量名称
	field                    string   // 成员变量，实际为 typeAlias.fieldName
	NotAllowEmptyValue       bool     // 是否允许为空，表示json字符传是否可为空null
}

// BuildTypes gen types to string
func BuildTypes(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name()+" generate error")
		}

	}

	return builder.String(), nil
}

func genTypes(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	typeBuilder := typeBuilder{
		types: api.Types,
	}
	err := typeBuilder.getTypesContent()
	if err != nil {
		return err
	}

	typeFilename, err := format.FileNamingFormat(cfg.NamingFormat, typesFile)
	if err != nil {
		return err
	}

	typeFilename = typeFilename + ".go"
	filename := path.Join(dir, typesDir, typeFilename)
	os.Remove(filename)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          typesDir,
		filename:        typeFilename,
		templateName:    "typesTemplate",
		category:        category,
		templateFile:    typesTemplateFile,
		builtinTemplate: typesTemplate,
		data: map[string]interface{}{
			"types":              typeBuilder.typesContent,
			"importPackages":     strings.Join(typeBuilder.importPackagesContents, "\n"),
			"containsTime":       false,
			"vars":               strings.Join(typeBuilder.globalVarsContents, "\n"),
			"typeFuncValidators": strings.Join(typeBuilder.typeFuncValidatorsContents, "\n\n"),
		},
	})
}

func writeType(writer io.Writer, tp spec.Type) error {
	structType, ok := tp.(spec.DefineStruct)
	if !ok {
		return fmt.Errorf("unsupport struct type: %s", tp.Name())
	}

	fmt.Fprintf(writer, "type %s struct {\n", util.Title(tp.Name()))
	for _, member := range structType.Members {
		if member.IsInline {
			if _, err := fmt.Fprintf(writer, "%s\n", strings.Title(member.Type.Name())); err != nil {
				return err
			}

			continue
		}

		if err := writeProperty(writer, member.Name, member.Tag, member.GetComment(), member.Type, 1, member.IsRequired); err != nil {
			return err
		}
	}
	fmt.Fprintf(writer, "}")
	return nil
}

// getTypesContent获取生成结构体所需要的数据
func (g *typeBuilder) getTypesContent() error {
	typesContent, err := BuildTypes(g.types)
	if err != nil {
		return err
	}
	g.typesContent = typesContent
	return g.getValidatorContents()
}

// getValidatorContents获取生成结构体所需要的数据
func (g *typeBuilder) getValidatorContents() error {
	for _, tp := range g.types {
		structType, ok := tp.(spec.DefineStruct)
		if !ok {
			return fmt.Errorf("unspport struct type: %s", tp.Name())
		}
		typeFuncBuilder := typeValidateFuncBuilder{
			typeGen: g,
		}
		err := typeFuncBuilder.getValidatorContent(structType)
		if err != nil {
			return err
		}
	}
	return nil
}

// getValidatorContent获取生成结构体所需要的数据
func (g *typeValidateFuncBuilder) getValidatorContent(structType spec.DefineStruct) error {
	if err := g.getFieldValidateContents(structType); err != nil {
		return err
	}
	g.getNullValidateContents(structType)

	localVarsContentMap := make(map[string]struct{}, 0)
	noRepeatLocalVarsContents := make([]string, 0)
	for _, localVarsContent := range g.localVarsContents {
		if _, ok := localVarsContentMap[localVarsContent]; !ok {
			localVarsContentMap[localVarsContent] = struct{}{}
			noRepeatLocalVarsContents = append(noRepeatLocalVarsContents, localVarsContent)
		}
	}
	data := map[string]interface{}{
		"type":                  util.Title(structType.Name()),
		"localVarContents":      strings.Join(noRepeatLocalVarsContents, "\n"),
		"fieldNullContents":     strings.Join(g.propertyNullContents, "\n"),
		"fieldValidateContents": strings.Join(g.fieldValidateContents, "\n"),
	}
	funcBuffer := &bytes.Buffer{}
	err := typeFuncValidatorTpl.Execute(funcBuffer, data)
	if err != nil {
		return err
	}
	g.typeGen.typeFuncValidatorsContents = append(g.typeGen.typeFuncValidatorsContents, funcBuffer.String())
	return nil
}

// getNullValidateContents 获取成员变量空校验内容
func (g *typeValidateFuncBuilder) getNullValidateContents(structType spec.DefineStruct) {
	var propertyNullValidateTpl = `
if %s == nil {
	return fmt.Errorf("%s is null")
}
`
	for _, member := range structType.Members {
		if !member.IsRequired {
			continue
		}
		propertyName := fmt.Sprintf("%s.%s", _validateFuncTypeAlias, util.Title(member.Name))
		nullValidateContent := fmt.Sprintf(propertyNullValidateTpl, propertyName, propertyName)
		g.propertyNullContents = append(g.propertyNullContents, nullValidateContent)

		g.typeGen.importPackagesContents = append(g.typeGen.importPackagesContents, _packageFmt)
	}
}

// getNullValidateContents 获取结构所有的成员变量校验内容
func (g *typeValidateFuncBuilder) getFieldValidateContents(structType spec.DefineStruct) error {
	for _, member := range structType.Members {
		field := fmt.Sprintf("%s.%s", _validateFuncTypeAlias, util.Title(member.Name))
		field = requiredTypeName(field, member.IsRequired, member.Type.Nullable())
		typeFieldValidateBuilder := typeFieldValidateBuilder{
			typeValidateFuncBuilder: g,
			typeAlias:               _validateFuncTypeAlias,
			isRequire:               member.IsRequired,
			rootTypeName:            structType.Name(),
			fieldName:               util.Title(member.Name),
			field:                   field,
			NotAllowEmptyValue:      member.NotAllowEmptyValue,
		}
		fieldValidateContent, err := typeFieldValidateBuilder.getFieldValidateContent(member.Type)
		if err != nil {
			return err
		}
		if fieldValidateContent == "" {
			continue
		}
		g.fieldValidateContents = append(g.fieldValidateContents, fieldValidateContent)
	}
	return nil
}

// getNullValidateContent 获取成员变量校验内容
func (g *typeFieldValidateBuilder) getFieldValidateContent(specType spec.Type) (string, error) {
	if defineStructSpecType, ok := specType.(spec.DefineStruct); ok {
		return g.getObjectTypeValidator(defineStructSpecType)
	}
	if primitiveType, ok := specType.(spec.PrimitiveType); ok {
		return g.getPrimitiveTypeValidator(primitiveType)
	}
	if primitiveType, ok := specType.(spec.ArrayType); ok {
		return g.getArrayValidator(primitiveType)
	}
	if mapType, ok := specType.(spec.MapType); ok {
		return g.getMapValidator(mapType)
	}
	return "", nil
}

// getObjectTypeValidator 获取基本对象型成员变量校验内容
func (g *typeFieldValidateBuilder) getObjectTypeValidator(structType spec.DefineStruct) (string, error) {
	data := make(map[string]interface{}, 0)
	data["type"] = "object"
	data["property"] = strings.Replace(g.field, "*", "", 1)
	data["propertyName"] = g.fieldName
	return g.getFieldValidContentWithTpl(data)
}

// getPrimitiveTypeValidatorObject 获取基本类型成员变量校验内容
func (g *typeFieldValidateBuilder) getPrimitiveTypeValidator(primitiveSpecType spec.PrimitiveType) (string, error) {
	numberTypeMap := map[string]struct{}{
		"int64": struct{}{}, "int32": struct{}{}, "int16": struct{}{}, "int8": struct{}{},
		"uint64": struct{}{}, "uint32": struct{}{}, "uint16": struct{}{}, "uint8": struct{}{},
		"float64": struct{}{}, "float32": struct{}{},
	}
	if _, ok := numberTypeMap[primitiveSpecType.RawName]; ok {
		return g.getPrimitiveTypeValidatorNumber(primitiveSpecType)
	}
	if primitiveSpecType.RawName == "string" {
		return g.getPrimitiveTypeValidatorString(primitiveSpecType)
	}
	return "", nil
}

// getFieldValidContentWithTpl 根据模板，填充数据返回具体的校验内容
func (g *typeFieldValidateBuilder) getFieldValidContentWithTpl(data map[string]interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	if err := typeFieldValidatorTpl.Execute(buffer, data); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// getFieldValidContentWithTpl 根据模板，填充数据返回具体的校验内容
func (g *typeFieldValidateBuilder) getPrimitiveTypeValidatorString(primitiveSpecType spec.PrimitiveType) (string, error) {
	data := make(map[string]interface{})
	if primitiveSpecType.MinLength != 0 {
		data["minLength"] = primitiveSpecType.MinLength
	}
	if primitiveSpecType.MaxLength != nil {
		data["shouldValidMaxLength"] = true
		data["maxLength"] = primitiveSpecType.MaxLength
	}
	if primitiveSpecType.Pattern != "" {
		newParentFieldNames := make([]string, 0)
		for _, v := range newParentFieldNames {
			newParentFieldNames = append(newParentFieldNames, util.Untitle(v))
		}
		regexPrefixNames := append([]string{util.Untitle(g.rootTypeName)}, newParentFieldNames...)
		regexName := fmt.Sprintf("%s_%s_regex", strings.Join(regexPrefixNames, "_"), util.Untitle(g.fieldName))
		regexName = util.Untitle(regexName)
		data["reg"] = regexName
		//depend on this var regexp, and should define in var
		varContent := fmt.Sprintf("%s = regexp.MustCompile(\"%s\")", regexName, primitiveSpecType.Pattern)
		g.typeGen.globalVarsContents = append(g.typeGen.globalVarsContents, varContent)

		g.typeGen.importPackagesContents = append(g.typeGen.importPackagesContents, _packageRegexp)
	}
	if g.NotAllowEmptyValue {
		data["notAllowEmptyValue"] = true
	}

	if primitiveSpecType.Enum != nil && len(primitiveSpecType.Enum) != 0 {
		enumStrs := make([]string, 0)
		for _, item := range primitiveSpecType.Enum {
			itemStr, ok := item.(string)
			if !ok {
				return "", fmt.Errorf("field: %s error enum: %v type, ", g.fieldName, item)
			}
			enumStrs = append(enumStrs, fmt.Sprintf(`"%s"`, itemStr))
		}
		enumDefine := fmt.Sprintf(" enumArray = []interface{}{%s}", strings.Join(enumStrs, ", "))
		data["enumValidate"] = true
		data["enumDefine"] = enumDefine
		g.localVarsContents = append(g.localVarsContents, "var enumArray []interface{}")
		g.localVarsContents = append(g.localVarsContents, "var enumExist bool")
	}
	if len(data) == 0 {
		return "", nil
	}

	g.typeGen.importPackagesContents = append(g.typeGen.importPackagesContents, _packageFmt)
	data["type"] = "string"
	data["property"] = g.field
	data["propertyName"] = g.fieldName

	return g.getFieldValidContentWithTpl(data)
}

// getPrimitiveTypeValidatorNumber 获取number型校验内容
func (g *typeFieldValidateBuilder) getPrimitiveTypeValidatorNumber(primitiveSpecType spec.PrimitiveType) (string, error) {
	data := make(map[string]interface{})
	if primitiveSpecType.Min != nil {
		data["shouldValidMin"] = true
		data["min"] = *primitiveSpecType.Min
	}
	if primitiveSpecType.Max != nil {
		data["shouldValidMax"] = true
		data["max"] = *primitiveSpecType.Max
	}
	if primitiveSpecType.ExclusiveMin {
		data["exclusiveMin"] = true
	}
	if primitiveSpecType.ExclusiveMax {
		data["exclusiveMax"] = true
	}
	if primitiveSpecType.MultipleOf != nil {
		data["multipleOf"] = *primitiveSpecType.MultipleOf
	}

	if g.NotAllowEmptyValue {
		data["notAllowEmptyValue"] = true
	}

	if primitiveSpecType.Enum != nil && len(primitiveSpecType.Enum) != 0 {
		enumNumber := make([]string, 0)
		for _, item := range primitiveSpecType.Enum {
			_, ok := item.(float64)
			if !ok {
				return "", fmt.Errorf("field: %s error enum: %v type, ", g.fieldName, item)
			}
			enumNumber = append(enumNumber, fmt.Sprint(item))
		}
		enumDefine := fmt.Sprintf(" enumArray = []interface{}{%s}", strings.Join(enumNumber, ", "))
		data["enumValidate"] = true
		data["enumDefine"] = enumDefine
		g.localVarsContents = append(g.localVarsContents, "var enumArray []interface{}")
		g.localVarsContents = append(g.localVarsContents, "var enumExist bool")
	}

	if len(data) == 0 {
		return "", nil
	}
	g.typeGen.importPackagesContents = append(g.typeGen.importPackagesContents, _packageFmt)

	data["type"] = "number"
	data["property"] = g.field
	data["propertyName"] = g.fieldName
	return g.getFieldValidContentWithTpl(data)
}

// getArrayValidator 获取array型校验内容
func (g *typeFieldValidateBuilder) getArrayValidator(arraySpecType spec.ArrayType) (string, error) {
	data := make(map[string]interface{})
	if arraySpecType.MinItems > 0 {
		data["minItems"] = arraySpecType.MinItems
	}
	if arraySpecType.MaxItems != nil {
		data["maxItems"] = *arraySpecType.MaxItems
	}
	if arraySpecType.UniqueItems {
		g.localVarsContents = append(g.localVarsContents, "var arrayRepeatMap map[interface{}]struct{}")
		data["shouldValidUniqueItems"] = true
	}
	if g.NotAllowEmptyValue {
		data["notAllowEmptyValue"] = true
	}
	arrayItemName := fmt.Sprintf("%sItem", util.Untitle(g.fieldName))
	typeFieldValidateGen := typeFieldValidateBuilder{
		typeValidateFuncBuilder: g.typeValidateFuncBuilder,
		typeAlias:               _validateFuncTypeAlias,
		rootTypeName:            g.rootTypeName,
		fieldName:               arrayItemName,
		field:                   arrayItemName,
		parentFieldNames:        append(g.parentFieldNames, g.fieldName),
	}
	fieldValidateContent, err := typeFieldValidateGen.getFieldValidateContent(arraySpecType.Value)
	if err != nil {
		return "", err
	}
	if fieldValidateContent != "" {
		data["itemValidateContent"] = fieldValidateContent
	}
	if len(data) == 0 {
		return "", nil
	}

	g.typeGen.importPackagesContents = append(g.typeGen.importPackagesContents, _packageFmt)

	data["type"] = "array"
	data["property"] = g.field
	data["propertyName"] = g.fieldName

	data["item"] = arrayItemName

	return g.getFieldValidContentWithTpl(data)
}

// getArrayValidator 获取array型校验内容
func (g *typeFieldValidateBuilder) getMapValidator(mapSpecType spec.MapType) (string, error) {
	data := make(map[string]interface{})
	valueName := fmt.Sprintf("%sValue", util.Untitle(g.fieldName))
	typeFieldValidateGen := typeFieldValidateBuilder{
		typeValidateFuncBuilder: g.typeValidateFuncBuilder,
		typeAlias:               _validateFuncTypeAlias,
		rootTypeName:            g.rootTypeName,
		fieldName:               valueName,
		field:                   valueName,
		parentFieldNames:        append(g.parentFieldNames, g.fieldName),
	}
	fieldValidateContent, err := typeFieldValidateGen.getFieldValidateContent(mapSpecType.Value)
	if err != nil {
		return "", err
	}
	if fieldValidateContent != "" {
		data["valueValidateContent"] = fieldValidateContent
	}
	if len(data) == 0 {
		return "", nil
	}

	g.typeGen.importPackagesContents = append(g.typeGen.importPackagesContents, _packageFmt)

	data["type"] = "map"
	data["property"] = g.field
	data["propertyName"] = g.fieldName

	data["value"] = valueName

	return g.getFieldValidContentWithTpl(data)
}
