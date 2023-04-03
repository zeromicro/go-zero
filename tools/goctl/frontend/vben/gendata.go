// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vben

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	util2 "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genData(g *GenContext) error {
	var basicData, searchFormData, formData strings.Builder
	var useBaseInfo bool
	var statusBasicColumnData, statusFormColumnData string
	// generate basic and search form data
	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			for _, val := range specData.Members {
				if val.Name == "" {
					tmpType, _ := val.Type.(spec.DefineStruct)
					if tmpType.Name() == "BaseIDInfo" || tmpType.Name() == "BaseUUIDInfo" {
						useBaseInfo = true
					}
				} else if val.Name == "Status" {
					statusRenderData := bytes.NewBufferString("")
					protoTmpl, _ := template.New("proto").Parse(statusRenderTpl)
					_ = protoTmpl.Execute(statusRenderData, map[string]any{
						"modelName": strings.TrimSuffix(specData.RawName, "Info"),
					})
					statusBasicColumnData = statusRenderData.String()

					statusFormColumnData = fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    component: 'RadioButtonGroup',\n"+
						"    defaultValue: 1,\n    componentProps: {\n      options: [\n        { label: t('common.on'), value: 1 },\n    "+
						"    { label: t('common.off'), value: 2 },\n      ],\n    },\n  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
					)
				} else {
					basicData.WriteString(fmt.Sprintf("\n  {\n    title: t('%s'),\n    dataIndex: '%s',\n    width: 100,\n  },",
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)), strcase.ToLowerCamel(val.Name)))

					formData.WriteString(fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    %s\n    required: true,\n%s  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
						getComponent(val.Type.Name()),
						GetRules(val),
					))
				}
			}

			// put here in order to put status in the end
			if g.HasStatus {
				basicData.WriteString(statusBasicColumnData)
				formData.WriteString(statusFormColumnData)
			}
		}

		if v.Name() == fmt.Sprintf("%sListReq", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get field")
			}

			for _, val := range specData.Members {
				if val.Name != "" {
					searchFormData.WriteString(fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    component: 'Input',\n    colProps: { span: 8 },\n%s  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "ListReq")),
							strcase.ToLowerCamel(val.Name)), GetRules(val),
					))
				}
			}
		}
	}

	if err := util.With("dataTpl").Parse(dataTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"folderName":          g.FolderName,
		"basicData":           basicData.String(),
		"searchFormData":      searchFormData.String(),
		"formData":            formData.String(),
		"useBaseInfo":         useBaseInfo,
		"useUUID":             g.UseUUID,
		"hasStatus":           g.HasStatus,
	},
		filepath.Join(g.ViewDir, fmt.Sprintf("%s.data.ts", strcase.ToLowerCamel(g.ModelName))), g.Overwrite); err != nil {
		return err
	}
	return nil
}

func getComponent(dataType string) string {
	switch dataType {
	case "string":
		return "component: 'Input',"
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "int64", "uint32", "uint64", "float32", "float64":
		return "component: 'InputNumber',"
	case "bool":
		return "component: 'RadioButtonGroup',\n    defaultValue: false,\n    componentProps: {\n      options: [\n        { label: t('common.on'), value: false },\n        { label: t('common.off'), value: true },\n      ],\n    },"
	default:
		return "component: 'Input',"
	}
}

// GetRules returns the rules from tag.
func GetRules(t spec.Member) string {
	validatorString := util2.ExtractValidateString(t.Tag)
	if validatorString == "" {
		return ""
	}

	rules, err := ConvertTagToRules(validatorString)
	if err != nil {
		return ""
	}

	switch GetRuleType(t.Type.Name()) {
	case "string":
		return fmt.Sprintf("    rules: [{ %s }],\n", strings.Join(rules, ", "))
	case "number":
		return fmt.Sprintf("    rules: [{ type: 'number', %s }],\n", strings.Join(rules, ", "))
	case "float":
		return fmt.Sprintf("    rules: [{ type: 'float', %s }],\n", strings.Join(rules, ", "))
	default:
		return ""
	}
}

// GetRuleType returns the rule type from go type.
func GetRuleType(t string) string {
	switch t {
	case "string":
		return "string"
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "int64", "uint32", "uint64":
		return "number"
	case "float32", "float64":
		return "float"
	default:
		return "string"
	}
}

// ConvertTagToRules converts validator tag to rules.
func ConvertTagToRules(tagString string) ([]string, error) {
	var result []string
	vals := strings.Split(tagString, ",")
	for _, v := range vals {
		if strings.Contains(v, "min") || strings.Contains(v, "max") {
			result = append(result, strings.Replace(v, "=", ": ", -1))
		}

		if strings.Contains(v, "len") {
			tagSplit := strings.Split(v, "=")
			_, tagNum := tagSplit[0], tagSplit[1]
			result = append(result, fmt.Sprintf("len: %s", tagNum))
		}

		if strings.Contains(v, "gt") || strings.Contains(v, "gte") ||
			strings.Contains(v, "lt") || strings.Contains(v, "lte") {
			tagSplit := strings.Split(v, "=")
			tag, tagNum := tagSplit[0], tagSplit[1]
			if strings.Contains(tagNum, ".") {
				bitSize := len(tagNum) - strings.Index(tagNum, ".") - 1
				n, err := strconv.ParseFloat(tagNum, bitSize)
				if err != nil {
					return nil, errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					result = append(result, fmt.Sprintf("min: %.*f", bitSize, n))
				case "gt":
					result = append(result, fmt.Sprintf("min: %.*f", bitSize, n+1/math.Pow(10, float64(bitSize))))
				case "lte":
					result = append(result, fmt.Sprintf("max: %.*f", bitSize, n))
				case "lt":
					result = append(result, fmt.Sprintf("max: %.*f", bitSize, n-1/math.Pow(10, float64(bitSize))))
				}
			} else {
				n, err := strconv.Atoi(tagNum)
				if err != nil {
					return nil, errors.New("failed to convert the number in validate tag")
				}

				switch tag {
				case "gte":
					result = append(result, fmt.Sprintf("min: %d", n))
				case "gt":
					result = append(result, fmt.Sprintf("min: %d", n))
				case "lte":
					result = append(result, fmt.Sprintf("max: %d", n))
				case "lt":
					result = append(result, fmt.Sprintf("max: %d", n))
				}
			}

		}
	}
	return result, nil
}
