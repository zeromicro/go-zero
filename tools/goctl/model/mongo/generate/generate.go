package generate

import (
	"errors"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/model/mongo/template"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

// Context defines the model generation data what they needs
type Context struct {
	Types  []string
	Cache  bool
	Easy   bool
	Output string
	Cfg    *config.Config
}

// Do executes model template and output the result into the specified file path
func Do(ctx *Context) error {
	if ctx.Cfg == nil {
		return errors.New("missing config")
	}

	if err := generateTypes(ctx); err != nil {
		return err
	}

	if err := generateModel(ctx); err != nil {
		return err
	}

	if err := generateCustomModel(ctx); err != nil {
		return err
	}

	return generateError(ctx)
}

func generateModel(ctx *Context) error {
	for _, t := range ctx.Types {
		fn, err := format.FileNamingFormat(ctx.Cfg.NamingFormat, t+"_model_gen")
		if err != nil {
			return err
		}

		text, err := pathx.LoadTemplate(category, modelTemplateFile, template.ModelText)
		if err != nil {
			return err
		}

		output := filepath.Join(ctx.Output, fn+".go")
		if err = util.With("model").Parse(text).GoFmt(true).SaveTo(map[string]any{
			"Type":      stringx.From(t).Title(),
			"lowerType": stringx.From(t).Untitle(),
			"Cache":     ctx.Cache,
			"version":   version.BuildVersion,
		}, output, true); err != nil {
			return err
		}
	}

	return nil
}

func generateCustomModel(ctx *Context) error {
	for _, t := range ctx.Types {
		fn, err := format.FileNamingFormat(ctx.Cfg.NamingFormat, t+"_model")
		if err != nil {
			return err
		}

		text, err := pathx.LoadTemplate(category, modelCustomTemplateFile, template.ModelCustomText)
		if err != nil {
			return err
		}

		output := filepath.Join(ctx.Output, fn+".go")
		err = util.With("model").Parse(text).GoFmt(true).SaveTo(map[string]any{
			"Type":      stringx.From(t).Title(),
			"lowerType": stringx.From(t).Untitle(),
			"snakeType": stringx.From(t).ToSnake(),
			"Cache":     ctx.Cache,
			"Easy":      ctx.Easy,
		}, output, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTypes(ctx *Context) error {
	for _, t := range ctx.Types {
		fn, err := format.FileNamingFormat(ctx.Cfg.NamingFormat, t+"_types")
		if err != nil {
			return err
		}

		text, err := pathx.LoadTemplate(category, modelTypesTemplateFile, template.ModelTypesText)
		if err != nil {
			return err
		}

		output := filepath.Join(ctx.Output, fn+".go")
		if err = util.With("model").Parse(text).GoFmt(true).SaveTo(map[string]any{
			"Type": stringx.From(t).Title(),
		}, output, false); err != nil {
			return err
		}
	}

	return nil
}

func generateError(ctx *Context) error {
	text, err := pathx.LoadTemplate(category, errTemplateFile, template.Error)
	if err != nil {
		return err
	}

	output := filepath.Join(ctx.Output, "error.go")

	return util.With("error").Parse(text).GoFmt(true).SaveTo(ctx, output, false)
}
