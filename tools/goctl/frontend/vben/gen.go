package vben

import (
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
	}

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
