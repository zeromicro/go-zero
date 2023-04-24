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
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func genLocale(g *GenContext) error {
	var localeEnData, localeZhData strings.Builder
	var enLocaleFileName, zhLocaleFileName string
	enLocaleFileName = filepath.Join(g.LocaleDir, "en", fmt.Sprintf("%s.ts", strings.ToLower(g.FolderName)))
	zhLocaleFileName = filepath.Join(g.LocaleDir, "zh-CN", fmt.Sprintf("%s.ts", strings.ToLower(g.FolderName)))

	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			localeEnData.WriteString(fmt.Sprintf("  %s: {\n", strcase.ToLowerCamel(g.ModelName)))
			localeZhData.WriteString(fmt.Sprintf("  %s: {\n", strcase.ToLowerCamel(g.ModelName)))

			for _, val := range specData.Members {
				if val.Name != "" {
					localeEnData.WriteString(fmt.Sprintf("    %s: '%s',\n",
						strcase.ToLowerCamel(val.Name), strcase.ToCamel(val.Name)))

					localeZhData.WriteString(fmt.Sprintf("    %s: '%s',\n",
						strcase.ToLowerCamel(val.Name), strcase.ToCamel(val.Name)))
				}
			}

			localeEnData.WriteString(fmt.Sprintf("    add%s: 'Add %s',\n", g.ModelName, g.ModelName))
			localeEnData.WriteString(fmt.Sprintf("    edit%s: 'Edit %s',\n", g.ModelName, g.ModelName))
			localeEnData.WriteString(fmt.Sprintf("    %sList: '%s List',\n", strcase.ToLowerCamel(g.ModelName), g.ModelName))
			localeEnData.WriteString("  },\n")

			localeZhData.WriteString(fmt.Sprintf("    add%s: '添加 %s',\n", g.ModelName, g.ModelName))
			localeZhData.WriteString(fmt.Sprintf("    edit%s: '编辑 %s',\n", g.ModelName, g.ModelName))
			localeZhData.WriteString(fmt.Sprintf("    %sList: '%s 列表',\n", strcase.ToLowerCamel(g.ModelName), g.ModelName))
			localeZhData.WriteString("  },\n")
		}
	}

	if !pathx.FileExists(enLocaleFileName) || g.Overwrite {
		if err := util.With("localeTpl").Parse(localeTpl).SaveTo(map[string]any{
			"localeData": localeEnData.String(),
		},
			enLocaleFileName, g.Overwrite); err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(enLocaleFileName)
		if err != nil {
			return err
		}

		data := string(file)

		if !strings.Contains(data, strings.ToLower(g.ModelName)+":") {
			data = data[:len(data)-3] + localeEnData.String() + data[len(data)-3:]
		} else if g.Overwrite {
			begin, end := FindBeginEndOfLocaleField(data, strings.ToLower(g.ModelName))
			data = data[:begin-2] + localeEnData.String() + data[end+1:]
		} else {

		}

		err = os.WriteFile(enLocaleFileName, []byte(data), os.ModePerm)
		if err != nil {
			return err
		}
	}

	if !pathx.FileExists(zhLocaleFileName) {
		if err := util.With("localeTpl").Parse(localeTpl).SaveTo(map[string]any{
			"localeData": localeZhData.String(),
		},
			zhLocaleFileName, false); err != nil {
			return err
		}
	} else {
		file, err := os.ReadFile(zhLocaleFileName)
		if err != nil {
			return err
		}

		data := string(file)

		if !strings.Contains(data, strings.ToLower(g.ModelName)+":") {
			data = data[:len(data)-3] + localeZhData.String() + data[len(data)-3:]
		} else if g.Overwrite {
			begin, end := FindBeginEndOfLocaleField(data, strings.ToLower(g.ModelName))
			data = data[:begin-2] + localeZhData.String() + data[end+1:]
		}

		err = os.WriteFile(zhLocaleFileName, []byte(data), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
