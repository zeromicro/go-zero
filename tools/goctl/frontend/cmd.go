package frontend

import (
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/frontend/vben"
)

var (
	// Cmd describes an api command.
	Cmd = &cobra.Command{
		Use:   "frontend",
		Short: "Generate frontend related files",
	}

	vbenCmd = &cobra.Command{
		Use:   "vben",
		Short: "Generate frontend related files",
		RunE:  vben.GenCRUDLogic,
	}
)

func init() {
	vbenCmd.Flags().StringVar(&vben.VarStringOutput, "o", "./", "The output directory, it should be the root "+
		"directory of simple admin backend ui. ")
	vbenCmd.Flags().StringVar(&vben.VarStringApiFile, "apiFile", "", "The absolute path of api file.")
	vbenCmd.Flags().StringVar(&vben.VarStringFolderName, "folderName", "sys", "The folder name to generate file"+
		"in different directory. e.g. file folder in simple admin backend ui which is to store file manager service files. ")
	vbenCmd.Flags().StringVar(&vben.VarStringApiPrefix, "prefix", "sys-api", "The request prefix for proxy. e.g. sys-api ")
	vbenCmd.Flags().StringVar(&vben.VarStringModelName, "modelName", "", "The model name. e.g. Example ")

	Cmd.AddCommand(vbenCmd)
}
