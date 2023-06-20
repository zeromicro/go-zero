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

package makefile

import (
	"bytes"
	_ "embed"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suyuan32/knife/core/io/filex"

	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

//go:embed makefile.tpl
var makefileTpl string

var (
	VarStringServiceName string
	VarStringStyle       string
	VarStringDir         string
	VarStringServiceType string
	VarBoolI18n          bool
	VarBoolEnt           bool
)

type GenContext struct {
	ServiceName string
	Style       string
	IsSingle    bool
	IsApi       bool
	IsRpc       bool
	UseI18n     bool
	UseEnt      bool
	TargetPath  string
	EntFeature  string
}

func Gen(_ *cobra.Command, _ []string) (err error) {
	ctx := GenContext{}

	var absPath string

	if VarStringDir != "" {
		absPath, err = filepath.Abs(VarStringDir)
		if err != nil {
			return errors.Wrap(err, "dir not found")
		}
	} else {
		absPath, err = filepath.Abs(".")
		if err != nil {
			return errors.Wrap(err, "dir not found")
		}
	}

	filePath := filepath.Join(absPath, "Makefile")

	ctx.TargetPath = filePath
	ctx.Style = VarStringStyle
	ctx.ServiceName = VarStringServiceName
	ctx.UseEnt = VarBoolEnt
	ctx.UseI18n = VarBoolI18n

	switch VarStringServiceType {
	case "api":
		ctx.IsApi = true
	case "single":
		ctx.IsSingle = true
	case "rpc":
		ctx.IsRpc = true
	}

	if err := filex.Exist(ctx.TargetPath); err == nil {
		err = extractInfo(&ctx)
		if err != nil {
			return errors.Wrap(err, "failed to extract makefile info")
		}
	}

	err = DoGen(&ctx)

	return err
}

func DoGen(g *GenContext) error {
	serviceNameStyle, err := format.FileNamingFormat(g.Style, g.ServiceName)
	if err != nil {
		return errors.Wrap(err, "failed to format service name")
	}

	makefileData := bytes.NewBufferString("")
	makefileTmpl, _ := template.New("makefile").Parse(makefileTpl)
	_ = makefileTmpl.Execute(makefileData, map[string]any{
		"serviceName":      strcase.ToCamel(g.ServiceName),
		"useEnt":           g.UseEnt,
		"style":            g.Style,
		"useI18n":          g.UseI18n,
		"serviceNameStyle": serviceNameStyle,
		"serviceNameLower": strings.ToLower(g.ServiceName),
		"serviceNameSnake": strcase.ToSnake(g.ServiceName),
		"serviceNameDash":  strings.ReplaceAll(strcase.ToSnake(g.ServiceName), "_", "-"),
		"isApi":            g.IsApi,
		"isSingle":         g.IsSingle,
		"isRpc":            g.IsRpc,
		"entFeature":       g.EntFeature,
	})

	err = filex.RemoveIfExist(g.TargetPath)
	if err != nil {
		return err
	}

	err = filex.WriteFileString(g.TargetPath, makefileData.String(), filex.SuperReadWritePerm)
	if err != nil {
		return err
	}

	return nil
}
