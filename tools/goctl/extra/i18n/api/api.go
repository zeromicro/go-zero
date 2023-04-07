package api

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type GenContext struct {
	Target      string
	ModelName   string
	ModelNameZh string
	OutputDir   string
}

//go:embed zh.tpl
var zhTpl string

//go:embed en.tpl
var enTpl string

func GenApiI18n(g *GenContext) error {
	var zhString, enString strings.Builder
	var absPath string
	var err error

	if g.OutputDir != "" {
		absPath, err = filepath.Abs(g.OutputDir)
		if err != nil {
			return errors.Wrap(err, "failed to convert the output path")
		}
	}

	zhTemplate, err := template.New("i18n_zh").Parse(zhTpl)
	if err != nil {
		return errors.Wrap(err, "failed to create i18n api zh template")
	}

	err = zhTemplate.Execute(&zhString, map[string]any{
		"modelName":   g.ModelName,
		"modelNameZh": g.ModelNameZh,
	})
	if err != nil {
		return err
	}

	enTemplate, err := template.New("i18n_en").Parse(enTpl)
	if err != nil {
		return errors.Wrap(err, "failed to create i18n api zh template")
	}

	err = enTemplate.Execute(&enString, map[string]any{
		"modelName":   g.ModelName,
		"modelNameEn": strings.ToLower(strings.Replace(strcase.ToSnake(g.ModelName), "_", " ", -1)),
	})
	if err != nil {
		return err
	}

	if g.OutputDir != "" {
		zhPath := filepath.Join(absPath, "zh.json")
		if pathx.Exists(zhPath) {
			err = AppendToApiDesc(zhString.String(), zhPath)
			if err != nil {
				return err
			}
		}

		enPath := filepath.Join(absPath, "en.json")
		if pathx.Exists(enPath) {
			err = AppendToApiDesc(enString.String(), enPath)
			if err != nil {
				return err
			}
		}
	} else {
		console.Info(zhString.String() + "\n")
		console.Info(enString.String() + "\n")
	}
	return nil
}

func AppendToApiDesc(data, filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	originString := string(file)

	apiDescIndex := strings.Index(originString, "apiDesc")
	data = fmt.Sprintf("\n%s", data)

	var newString string
	offset := 12
	if apiDescIndex != -1 {
		newString = originString[:apiDescIndex+offset] + data + originString[apiDescIndex+offset:]
	}

	if newString != "" {
		err = os.WriteFile(filePath, []byte(newString), 0o666)
		if err != nil {
			return err
		}
	}

	return err
}
