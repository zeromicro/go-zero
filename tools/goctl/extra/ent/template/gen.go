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

package template

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/env"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"

	"github.com/suyuan32/knife/core/io/filex"
)

var (
	// VarStringDir is the ent directory
	VarStringDir string

	// VarStringAdd is the template name for adding
	VarStringAdd string

	// VarBoolUpdate describe whether to update the template
	VarBoolUpdate bool

	// VarBoolList describe whether to list all supported templates
	VarBoolList bool

	tplInfo table.Writer
)

func GenTemplate(_ *cobra.Command, _ []string) error {
	var entDir string
	var err error

	if VarBoolList {
		ListAllTemplate()
		return nil
	}

	if VarStringDir == "" {
		entDir, err = GetEntDir()
		if err != nil {
			return err
		}
	} else {
		entDir, err = filepath.Abs(VarStringDir)
		if err != nil {
			return err
		}
	}

	tmplDir := filepath.Join(entDir, "template")

	if VarBoolUpdate {
		files, err := filex.GetFilesPathFromDir(tmplDir, false)
		if err != nil {
			return err
		}

		for _, v := range files {
			fileName := filepath.Base(v)
			tpl := GetTmpl(fileName)
			if tpl == "" {
				return errors.New("failed to find the template")
			}

			err := filex.RemoveIfExist(v)
			if err != nil {
				return errors.Join(err, errors.New("failed to remove the original template"))
			}

			err = filex.WriteFileString(v, tpl, filex.SuperReadWritePerm)
			if err != nil {
				return err
			}
		}

		execx.Run("go get -u entgo.io/ent@latest", entDir)
	}

	if VarStringAdd != "" {
		tpl := GetTmpl(VarStringAdd)
		if tpl == "" {
			return errors.New("failed to find the template")
		}

		filePath := filepath.Join(tmplDir, VarStringAdd+".tmpl")

		if pathx.Exists(filePath) {
			return errors.New("the template already exists")
		}

		err := filex.WriteFileString(filePath, tpl, filex.SuperReadWritePerm)
		if err != nil {
			return err
		}
	}

	console.Success("Generating successfully")

	return nil
}

func GetEntDir() (string, error) {
	entDir, _ := filepath.Abs("./ent")

	if pathx.Exists(entDir) {
		return entDir, nil
	}

	entDir, _ = filepath.Abs("./rpc/ent")

	if pathx.Exists(entDir) {
		return entDir, nil
	}

	entDir, _ = filepath.Abs("./api/ent")

	if pathx.Exists(entDir) {
		return entDir, nil
	}

	return "", errors.New("failed to find the ent directory")
}

func GetTmpl(name string) string {
	switch name {
	case "not_empty_update.tmpl", "not_empty_update":
		return NotEmptyTmpl
	case "pagination.tmpl", "pagination":
		return PaginationTmpl
	}
	return ""
}

func ListAllTemplate() {
	type Info struct {
		Name  string
		Intro string
	}

	var data []Info
	tplInfo = table.NewWriter()

	if env.IsChinaEnv() {
		color.Green.Println("支持的模板:\n")
		tplInfo.AppendHeader(table.Row{"模板名称", "模板介绍"})
		data = []Info{
			{
				"not_empty_update",
				"Ent 非空模板，用于如果输入值为空则不更新",
			},
			{
				"pagination",
				"Ent 分页模板",
			},
		}
	} else {
		color.Green.Println("The templates supported:\n")
		tplInfo.AppendHeader(table.Row{"Name", "Introduction"})
		data = []Info{
			{
				"not_empty_update",
				"The template for updating the values when it is not empty",
			},
			{
				"pagination",
				"The template for paginating the data",
			},
		}
	}

	for _, v := range data {
		tplInfo.AppendRows([]table.Row{
			{
				v.Name,
				v.Intro,
			},
		})
	}

	fmt.Println(tplInfo.Render())

	if env.IsChinaEnv() {
		color.Green.Println("\n使用方法： goctls extra ent template -a not_empty_update -d ./ent ")
	} else {
		color.Green.Println("\nUsage: goctls extra ent template -a not_empty_update -d ./ent")
	}
}
