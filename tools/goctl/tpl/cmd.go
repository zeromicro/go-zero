package tpl

import "github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"

var (
	varStringHome     string
	varStringCategory string
	varStringName     string
	// Cmd describes a template command.
	Cmd       = cobrax.NewCommand("template")
	initCmd   = cobrax.NewCommand("init", cobrax.WithRunE(genTemplates))
	cleanCmd  = cobrax.NewCommand("clean", cobrax.WithRunE(cleanTemplates))
	updateCmd = cobrax.NewCommand("update", cobrax.WithRunE(updateTemplates))
	revertCmd = cobrax.NewCommand("revert", cobrax.WithRunE(revertTemplates))
)

func init() {
	initCmd.Flags().StringVar(&varStringHome, "home")
	cleanCmd.Flags().StringVar(&varStringHome, "home")
	updateCmd.Flags().StringVar(&varStringHome, "home")
	updateCmd.Flags().StringVarP(&varStringCategory, "category", "c")
	revertCmd.Flags().StringVar(&varStringHome, "home")
	revertCmd.Flags().StringVarP(&varStringCategory, "category", "c")
	revertCmd.Flags().StringVarP(&varStringName, "name", "n")

	Cmd.AddCommand(cleanCmd, initCmd, revertCmd, updateCmd)
}
