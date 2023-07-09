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

package info

import (
	"github.com/zeromicro/go-zero/tools/goctl/info/env"
	"github.com/zeromicro/go-zero/tools/goctl/info/port"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

var (
	// Cmd describes a docker command.
	Cmd = cobrax.NewCommand("info")

	EnvCmd = cobrax.NewCommand("env", cobrax.WithRunE(env.ShowEnv))

	PortCmd = cobrax.NewCommand("port", cobrax.WithRunE(port.ShowPort))
)

func init() {

	var (
		envCmdFlags = EnvCmd.Flags()
	)

	envCmdFlags.StringVarPWithDefaultValue(&env.ServiceName, "service_name", "s", "core")
	envCmdFlags.BoolVarP(&env.ShowList, "list", "l")

	Cmd.AddCommand(EnvCmd)
	Cmd.AddCommand(PortCmd)
}
