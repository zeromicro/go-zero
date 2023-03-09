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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genApi(g *GenContext) error {
	if err := util.With("apiTpl").Parse(apiTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameSpace":      strings.Replace(strcase.ToSnake(g.ModelName), "_", " ", -1),
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"modelNameSnake":      strcase.ToSnake(g.ModelName),
		"prefix":              g.Prefix,
		"useUUID":             g.UseUUID,
		"hasStatus":           g.HasStatus,
	},
		filepath.Join(g.ApiDir, fmt.Sprintf("%s.ts", strcase.ToLowerCamel(g.ModelName))), g.Overwrite); err != nil {
		return err
	}
	return nil
}
