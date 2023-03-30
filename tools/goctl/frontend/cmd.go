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

	vbenCmdFlags.StringVarWithDefaultValue(&vben.VarStringOutput, "o", "./")
	vbenCmdFlags.StringVar(&vben.VarStringApiFile, "api_file")
	vbenCmdFlags.StringVarWithDefaultValue(&vben.VarStringFolderName, "folder_name", "sys")
	vbenCmdFlags.StringVar(&vben.VarStringSubFolder, "sub_folder")
	vbenCmdFlags.StringVarWithDefaultValue(&vben.VarStringApiPrefix, "prefix", "sys-api")
	vbenCmdFlags.StringVar(&vben.VarStringModelName, "model_name")
	vbenCmdFlags.BoolVar(&vben.VarBoolOverwrite, "overwrite")

	Cmd.AddCommand(VbenCmd)
}
