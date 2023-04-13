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

package new

import (
	_ "embed"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringHome describes the goctl home.
	VarStringHome string
	// VarStringRemote describes the remote git repository.
	VarStringRemote string
	// VarStringBranch describes the git branch.
	VarStringBranch string
	// VarStringStyle describes the style of output files.
	VarStringStyle string
	// VarBoolErrorTranslate describes whether to translate error
	VarBoolErrorTranslate bool
	// VarBoolUseCasbin describe whether to use Casbin
	VarBoolUseCasbin bool
	// VarBoolUseI18n describe whether to use i18n
	VarBoolUseI18n bool
	// VarStringGoZeroVersion describe the version of Go Zero
	VarStringGoZeroVersion string
	// VarStringToolVersion describe the version of Simple Admin Tools
	VarStringToolVersion string
	// VarModuleName describe the module name
	VarModuleName string
	// VarIntServicePort describe the service port exposed
	VarIntServicePort int
	// VarBoolGitlab describes whether to use gitlab-ci
	VarBoolGitlab bool
	// VarBoolEnt describes whether to use ent in api
	VarBoolEnt bool
)

// CreateServiceCommand fast create service
func CreateServiceCommand(_ *cobra.Command, args []string) error {
	dirName := args[0]
	if len(VarStringStyle) == 0 {
		VarStringStyle = conf.DefaultFormat
	}
	if strings.Contains(dirName, "-") {
		return errors.New("api new command service name not support strikethrough, because this will used by function name")
	}

	abs, err := filepath.Abs(dirName)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(filepath.Join(abs, "desc"))
	if err != nil {
		return err
	}

	if len(VarStringRemote) > 0 {
		repo, _ := util.CloneIntoGitHome(VarStringRemote, VarStringBranch)
		if len(repo) > 0 {
			VarStringHome = repo
		}
	}

	if len(VarStringHome) > 0 {
		pathx.RegisterGoctlHome(VarStringHome)
	}

	apiFilePath := filepath.Join(abs, "desc", "all.api")

	text, err := pathx.LoadTemplate(category, apiTemplateFile, baseApiTmpl)
	if err != nil {
		return err
	}

	baseApiFile, err := os.Create(filepath.Join(abs, "desc", "base.api"))
	if err != nil {
		return err
	}
	defer baseApiFile.Close()

	t := template.Must(template.New("baseApiTemplate").Parse(text))
	if err := t.Execute(baseApiFile, map[string]string{
		"name": strcase.ToCamel(dirName),
	}); err != nil {
		return err
	}

	allApiFile, err := os.Create(filepath.Join(abs, "desc", "all.api"))
	if err != nil {
		return err
	}
	defer allApiFile.Close()

	allTpl := template.Must(template.New("allApiTemplate").Parse(allApiTmpl))
	if err := allTpl.Execute(allApiFile, map[string]string{
		"name": strcase.ToCamel(dirName),
	}); err != nil {
		return err
	}

	var moduleName string

	if VarModuleName != "" {
		moduleName = VarModuleName
	} else {
		moduleName = dirName
	}

	genCtx := &gogen.GenContext{
		GoZeroVersion: VarStringGoZeroVersion,
		ToolVersion:   VarStringToolVersion,
		UseCasbin:     VarBoolUseCasbin,
		UseI18n:       VarBoolUseI18n,
		TransErr:      VarBoolErrorTranslate,
		ModuleName:    moduleName,
		Port:          VarIntServicePort,
		UseGitlab:     VarBoolGitlab,
		UseMakefile:   true,
		UseDockerfile: true,
		UseEnt:        VarBoolEnt,
	}

	err = gogen.DoGenProject(apiFilePath, abs, VarStringStyle, genCtx)
	return err
}
