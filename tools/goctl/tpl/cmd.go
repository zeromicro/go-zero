package tpl

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

var (
	varStringHome     string
	varStringCategory string
	varStringName     string
	// Cmd describes a template command.
	Cmd = &cobra.Command{
		Use:   "template",
		Short: flags.Get("template.short"),
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: flags.Get("template.init.short"),
		RunE:  genTemplates,
	}

	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: flags.Get("template.clean.short"),
		RunE:  cleanTemplates,
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: flags.Get("template.update.short"),
		RunE:  updateTemplates,
	}

	revertCmd = &cobra.Command{
		Use:   "revert",
		Short: flags.Get("template.revert.short"),
		RunE:  revertTemplates,
	}
)

func init() {
	initCmd.Flags().StringVar(&varStringHome, "home", "", flags.Get("template.init.home"))
	cleanCmd.Flags().StringVar(&varStringHome, "home", "", flags.Get("template.clean.home"))
	updateCmd.Flags().StringVar(&varStringHome, "home", "", flags.Get("template.update.home"))
	updateCmd.Flags().StringVarP(&varStringCategory, "category", "c", "", flags.Get("template.update.category"))
	revertCmd.Flags().StringVar(&varStringHome, "home", "", flags.Get("template.revert.home"))
	revertCmd.Flags().StringVarP(&varStringCategory, "category", "c", "", flags.Get("template.revert.category"))
	revertCmd.Flags().StringVarP(&varStringName, "name", "n", "", flags.Get("template.revert.name"))

	Cmd.AddCommand(cleanCmd, initCmd, revertCmd, updateCmd)
}
