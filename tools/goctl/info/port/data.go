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

package port

import (
	"os"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/zeromicro/go-zero/tools/goctl/util/env"
)

var envInfo table.Writer
var lang = env.IsChinaEnv()

// portInfo show the port usage across the simple admin
func portInfo() string {
	color.Green.Println("PORT")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"端口", "服务"})

	} else {
		envInfo.AppendHeader(table.Row{"Port", "Service"})
	}
	envInfo.AppendRows([]table.Row{
		{9100, "core_api"},
		{9101, "core_rpc"},
		{9102, "file_api"},
		{9103, "member_api"},
		{9104, "member_rpc"},
		{9105, "job_rpc"},
		{9106, "mcms_rpc"},
	})
	return envInfo.Render()
}
