package frontend

import (
	"github.com/zeromicro/go-zero/tools/goctl/frontend/vben"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

var (
	// Cmd describes an api command.
	Cmd     = cobrax.NewCommand("frontend")
	VbenCmd = cobrax.NewCommand("vben", cobrax.WithRunE(vben.GenCRUDLogic))
)

func init() {
	vbenCmdFlags := VbenCmd.Flags()

	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringOutput, "output", "o", "./")
	vbenCmdFlags.StringVarP(&vben.VarStringApiFile, "api_file", "a")
	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringFolderName, "folder_name", "f", "sys")
	vbenCmdFlags.StringVarP(&vben.VarStringSubFolder, "sub_folder", "s")
	vbenCmdFlags.StringVarPWithDefaultValue(&vben.VarStringApiPrefix, "prefix", "p", "sys-api")
	vbenCmdFlags.StringVarP(&vben.VarStringModelName, "model_name", "m")
	vbenCmdFlags.BoolVarP(&vben.VarBoolOverwrite, "overwrite", "w")

	Cmd.AddCommand(VbenCmd)
}
