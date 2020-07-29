package gen

import (
	"fmt"
	"strings"
	"text/template"

	"zero/tools/goctl/api/spec"
	"zero/tools/goctl/api/util"
	"zero/tools/goctl/model/mongomodel/utils"
)

const (
	functionTypeGet  = "get"  //GetByField return single model
	functionTypeFind = "find" // findByField return many model
	functionTypeSet  = "set"  // SetField  only set specified field

	TagOperate = "o" //字段函数的tag
	TagComment = "c" //字段注释的tag
)

type (
	FunctionDesc struct {
		Type      string // get,find,set
		FieldName string // 字段名字 eg:Age
		FieldType string // 字段类型 eg: string,int64 等
	}
)

func GenMongoModel(goFilePath string, needCache bool) error {
	structs, imports, err := utils.ParseGoFile(goFilePath)
	if err != nil {
		return err
	}
	if len(structs) != 1 {
		return fmt.Errorf("only 1 struct should be provided")
	}
	structStr, err := genStructs(structs)
	if err != nil {
		return err
	}

	fp, err := util.ClearAndOpenFile(goFilePath)
	if err != nil {
		return err
	}
	defer fp.Close()

	var myTemplate string
	if needCache {
		myTemplate = cacheTemplate
	} else {
		myTemplate = noCacheTemplate
	}
	structName := getStructName(structs)
	functionList := getFunctionList(structs)

	for _, fun := range functionList {
		funTmp := genMethodTemplate(fun, needCache)
		if funTmp == "" {
			continue
		}
		myTemplate += "\n"
		myTemplate += funTmp
		myTemplate += "\n"
	}

	t := template.Must(template.New("mongoTemplate").Parse(myTemplate))
	return t.Execute(fp, map[string]string{
		"modelName":   structName,
		"importArray": getImports(imports, needCache),
		"modelFields": structStr,
	})
}

func getFunctionList(structs []utils.Struct) []FunctionDesc {
	var list []FunctionDesc
	for _, field := range structs[0].Fields {
		tagMap := parseTag(field.Tag)
		if fun, ok := tagMap[TagOperate]; ok {
			funList := strings.Split(fun, ",")
			for _, o := range funList {
				var f FunctionDesc
				f.FieldType = field.Type
				f.FieldName = field.Name
				f.Type = o
				list = append(list, f)
			}
		}
	}
	return list
}

func getStructName(structs []utils.Struct) string {
	for _, structItem := range structs {
		return structItem.Name
	}
	return ""
}

func genStructs(structs []utils.Struct) (string, error) {
	if len(structs) > 1 {
		return "", fmt.Errorf("input .go file must only one struct")
	}

	modelFields := `Id             bson.ObjectId ` + quotationMark + `bson:"_id" json:"id,omitempty"` + quotationMark + "\n\t"
	for _, structItem := range structs {
		for _, field := range structItem.Fields {
			modelFields += getFieldLine(field)
		}
	}
	modelFields += "\t" + `CreateTime time.Time ` + quotationMark + `json:"createTime,omitempty" bson:"createTime"` + quotationMark + "\n\t"
	modelFields += "\t" + `UpdateTime time.Time ` + quotationMark + `json:"updateTime,omitempty" bson:"updateTime"` + quotationMark
	return modelFields, nil
}

func getFieldLine(member spec.Member) string {
	if member.Name == "CreateTime" || member.Name == "UpdateTime" || member.Name == "Id" {
		return ""
	}
	jsonName := utils.UpperCamelToLower(member.Name)
	result := "\t" + member.Name + ` ` + member.Type + ` ` + quotationMark + `json:"` + jsonName + `,omitempty"` + ` bson:"` + jsonName + `"` + quotationMark
	tagMap := parseTag(member.Tag)
	if comment, ok := tagMap[TagComment]; ok {
		result += ` //` + comment + "\n\t"
	} else {
		result += "\n\t"
	}
	return result
}

// tag like `o:"find,get,update" c:"姓名"`
func parseTag(tag string) map[string]string {
	var result = make(map[string]string, 0)
	tags := strings.Split(tag, " ")
	for _, kv := range tags {
		temp := strings.Split(kv, ":")
		if len(temp) > 1 {
			key := strings.ReplaceAll(strings.ReplaceAll(temp[0], "\"", ""), "`", "")
			value := strings.ReplaceAll(strings.ReplaceAll(temp[1], "\"", ""), "`", "")
			result[key] = value
		}
	}
	return result
}

func getImports(imports []string, needCache bool) string {

	importStr := strings.Join(imports, "\n\t")
	importStr += "\"errors\"\n\t"
	importStr += "\"time\"\n\t"
	importStr += "\n\t\"zero/core/stores/cache\"\n\t"
	importStr += "\"zero/core/stores/mongoc\"\n\t"
	importStr += "\n\t\"github.com/globalsign/mgo/bson\""
	return importStr
}
