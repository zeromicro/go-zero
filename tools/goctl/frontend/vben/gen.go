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
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringOutput describes the output.
	VarStringOutput string
	// VarStringApiFile describes the api file path
	VarStringApiFile string
	// VarStringFolderName describes the folder name to output dir
	VarStringFolderName string
	// VarStringApiPrefix describes the request URL's prefix
	VarStringApiPrefix string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringSubFolder describes the sub folder name
	VarStringSubFolder string
	// VarBoolOverwrite describes whether to overwrite the files, it will overwrite all generated files.
	VarBoolOverwrite bool
)

type GenContext struct {
	ApiDir        string
	ModelDir      string
	ViewDir       string
	Prefix        string
	ModelName     string
	LocaleDir     string
	FolderName    string
	SubFolderName string
	ApiSpec       *spec.ApiSpec
	UseUUID       bool
	HasStatus     bool
	Overwrite     bool
}

func (g GenContext) Validate() error {
	if g.ApiDir == "" {
		return errors.New("please set the api file path via --api_file")
	} else if !strings.HasSuffix(g.ApiDir, "api") {
		return errors.New("please input correct api file path")
	}
	return nil
}

// GenCRUDLogic is used to generate CRUD file for simple admin backend UI
func GenCRUDLogic(_ *cobra.Command, _ []string) error {
	outputDir, err := filepath.Abs(VarStringOutput)
	if err != nil {
		return err
	}

	apiFile, err := parser.Parse(VarStringApiFile)
	if err != nil {
		return err
	}

	apiOutputDir := filepath.Join(outputDir, "src/api", VarStringFolderName)
	if err := pathx.MkdirIfNotExist(apiOutputDir); err != nil {
		return err
	}
	modelOutputDir := filepath.Join(outputDir, "src/api", VarStringFolderName, "model")
	if err := pathx.MkdirIfNotExist(modelOutputDir); err != nil {
		return err
	}
	viewOutputDir := filepath.Join(outputDir, "src/views", VarStringFolderName)
	if err := pathx.MkdirIfNotExist(viewOutputDir); err != nil {
		return err
	}
	if VarStringSubFolder != "" {
		viewOutputDir = filepath.Join(viewOutputDir, VarStringSubFolder)
		if err := pathx.MkdirIfNotExist(viewOutputDir); err != nil {
			return err
		}
	}
	localeDir := filepath.Join(outputDir, "src/locales/lang")

	var modelName string
	if VarStringModelName != "" {
		modelName = VarStringModelName
	} else {
		modelName = strcase.ToCamel(strings.TrimSuffix(filepath.Base(VarStringApiFile), ".api"))
	}

	genCtx := &GenContext{
		ApiDir:     apiOutputDir,
		ModelDir:   modelOutputDir,
		ViewDir:    viewOutputDir,
		Prefix:     VarStringApiPrefix,
		ModelName:  modelName,
		ApiSpec:    apiFile,
		LocaleDir:  localeDir,
		FolderName: VarStringFolderName,
		Overwrite:  VarBoolOverwrite,
	}

	err = genCtx.Validate()

	if err := genModel(genCtx); err != nil {
		return err
	}

	if err := genApi(genCtx); err != nil {
		return err
	}

	if err := genData(genCtx); err != nil {
		return err
	}

	if err := genLocale(genCtx); err != nil {
		return err
	}

	if err := genIndex(genCtx); err != nil {
		return err
	}

	if err := genDrawer(genCtx); err != nil {
		return err
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}
