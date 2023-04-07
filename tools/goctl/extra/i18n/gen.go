package i18n

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/extra/i18n/api"
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

func Gen(_ *cobra.Command, _ []string) error {
	err := Validate()
	if err != nil {
		return err
	}
	return DoGen()
}

func DoGen() error {
	switch VarStringTarget {
	case "api":
		ctx := &api.GenContext{
			Target:      VarStringTarget,
			ModelName:   VarStringModelName,
			ModelNameZh: VarStringModelNameZh,
			OutputDir:   VarStringOutputDir,
		}
		return api.GenApiI18n(ctx)
	}
	return errors.New("invalid target, try \"api\"")
}

func Validate() error {
	if VarStringTarget == "" {
		return errors.New("the target cannot be empty, use --target to set it")
	} else if VarStringModelName == "" {
		return errors.New("the model name cannot be empty, use --model_name to set it")
	}
	return nil
}
