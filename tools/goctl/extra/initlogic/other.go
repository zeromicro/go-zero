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

package initlogic

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

//go:embed other.tpl
var otherTpl string

func OtherGen(g *CoreGenContext) error {
	var otherString strings.Builder
	otherTemplate, err := template.New("init_other").Parse(otherTpl)
	if err != nil {
		return errors.Wrap(err, "failed to create other init template")
	}

	err = otherTemplate.Execute(&otherString, map[string]any{
		"modelName":      g.ModelName,
		"modelNameSnake": strcase.ToSnake(g.ModelName),
		"modelNameUpper": strings.ToUpper(g.ModelName),
	})
	if err != nil {
		return err
	}

	console.Info(otherString.String())

	return err
}
