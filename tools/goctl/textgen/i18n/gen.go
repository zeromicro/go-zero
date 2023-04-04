package i18n

import (
	"github.com/spf13/cobra"
)

var (
	// VarStringTarget describes the target.
	VarStringTarget string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringModelNameZh describes the model's Chinese translation
	VarStringModelNameZh string
	// VarStringOutputDir describes the output directory
	VarStringOutputDir string
)

type GenContext struct {
	Target      string
	ModelName   string
	ModelNameZh string
	OutputDir   string
}

func Gen(_ *cobra.Command, _ []string) error {
	ctx := &GenContext{
		Target:      VarStringTarget,
		ModelName:   VarStringModelName,
		ModelNameZh: VarStringModelNameZh,
		OutputDir:   VarStringOutputDir,
	}
	return DoGen(ctx)
}

func DoGen(ctx *GenContext) error {
	switch ctx.Target {
	case "api":
		return GenApiI18n(ctx)
	}
	return nil
}
