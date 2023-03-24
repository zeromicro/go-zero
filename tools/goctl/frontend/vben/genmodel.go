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
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genModel(g *GenContext) error {
	var infoData strings.Builder
	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			for _, val := range specData.Members {
				if val.Name == "" {
					tmpType, _ := val.Type.(spec.DefineStruct)
					if tmpType.Name() == "BaseIDInfo" {
						infoData.WriteString("  id: number;\n  createdAt?: number;\n  updatedAt?: number;\n")
					} else if tmpType.Name() == "BaseUUIDInfo" {
						infoData.WriteString("  id: string;\n  createdAt?: number;\n  updatedAt?: number;\n")
						g.UseUUID = true
					}
				} else {
					if val.Name == "Status" {
						g.HasStatus = true
					}

					infoData.WriteString(fmt.Sprintf("  %s?: %s;\n", strcase.ToLowerCamel(val.Name),
						ConvertGoTypeToTsType(val.Type.Name())))
				}
			}

		}
	}
	if err := util.With("modelTpl").Parse(modelTpl).SaveTo(map[string]any{
		"modelName": g.ModelName,
		"infoData":  infoData.String(),
	},
		filepath.Join(g.ModelDir, fmt.Sprintf("%sModel.ts", strcase.ToLowerCamel(g.ModelName))), g.Overwrite); err != nil {
		return err
	}
	return nil
}
