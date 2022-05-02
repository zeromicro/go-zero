package tpl

import (
	"github.com/spf13/cobra"
)

var (
	varStringHome     string
	varStringCategory string
	varStringName     string
	Cmd               = &cobra.Command{
		Use:   "template",
		Short: "template operation",
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize the all templates(force update)",
		RunE:  genTemplates,
	}

	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "clean the all cache templates",
		RunE:  cleanTemplates,
	}

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "update template of the target category to the latest",
		RunE:  updateTemplates,
	}

	revertCmd = &cobra.Command{
		Use:   "revert",
		Short: "revert the target template to the latest",
		RunE:  revertTemplates,
	}
)

func init() {
	initCmd.Flags().StringVar(&varStringHome, "home", "", "the goctl home path of the template")
	cleanCmd.Flags().StringVar(&varStringHome, "home", "", "the goctl home path of the template")
	updateCmd.Flags().StringVar(&varStringHome, "home", "", "the goctl home path of the template")
	updateCmd.Flags().StringVarP(&varStringCategory, "category", "c", "", "the category of template, enum [api,rpc,model,docker,kube]")
	revertCmd.Flags().StringVar(&varStringHome, "home", "", "the goctl home path of the template")
	revertCmd.Flags().StringVarP(&varStringCategory, "category", "c", "", "the category of template, enum [api,rpc,model,docker,kube]")
	revertCmd.Flags().StringVarP(&varStringName, "name", "n", "", "the target file name of template")

	Cmd.AddCommand(cleanCmd)
	Cmd.AddCommand(initCmd)
	Cmd.AddCommand(revertCmd)
	Cmd.AddCommand(updateCmd)
}
