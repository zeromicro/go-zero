// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
