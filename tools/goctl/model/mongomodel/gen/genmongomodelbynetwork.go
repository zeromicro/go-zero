package gen

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/model/mongomodel/utils"
)

func GenMongoModelByNetwork(input string, needCache bool) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("struct不能为空")
	}
	if strings.Index(strings.TrimSpace(input), "type") != 0 {
		input = "type " + input
	}

	if strings.Index(strings.TrimSpace(input), "package") != 0 {
		input = "package model\r\n" + input
	}

	structs, imports, err := utils.ParseGoFileByNetwork(input)
	if err != nil {
		return "", err
	}
	if len(structs) != 1 {
		return "", fmt.Errorf("only 1 struct should be provided")
	}
	structStr, err := genStructs(structs)
	if err != nil {
		return "", err
	}

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
	var result bytes.Buffer
	err = t.Execute(&result, map[string]string{
		"modelName":   structName,
		"importArray": getImports(imports, needCache),
		"modelFields": structStr,
	})
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
