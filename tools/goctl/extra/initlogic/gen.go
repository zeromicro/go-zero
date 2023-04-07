package initlogic

import (
	_ "embed"
	"errors"

	"github.com/spf13/cobra"
)

var (
	// VarStringTarget describes the target.
	VarStringTarget string
	// VarStringModelName describes the model name
	VarStringModelName string
	// VarStringOutputPath describes the output directory
	VarStringOutputPath string
)

func Gen(_ *cobra.Command, _ []string) error {
	err := Validate()
	if err != nil {
		return err
	}

	ctx := &CoreGenContext{
		Target:    VarStringTarget,
		ModelName: VarStringModelName,
		Output:    VarStringOutputPath,
	}

	return DoGen(ctx)
}

func DoGen(g *CoreGenContext) error {
	if g.Target == "core" {
		return GenCore(g)
	} else if g.Target == "other" {
		return OtherGen(g)
	}
	return errors.New("invalid target, try \"core\" or \"other\"")
}

func Validate() error {
	if VarStringTarget == "" {
		return errors.New("the target cannot be empty, use --target to set it")
	} else if VarStringModelName == "" {
		return errors.New("the model name cannot be empty, use --model_name to set it")
	}
	return nil
}
