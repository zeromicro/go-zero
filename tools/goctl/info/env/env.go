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

package env

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	ServiceName string
)

// ShowEnv is used to show the environment variable usages.
func ShowEnv(_ *cobra.Command, _ []string) error {
	if lang {
		color.Green.Println("Simple Admin的环境变量")
		color.Red.Println("注意： 环境变量的优先级大于配置文件")
	} else {
		color.Green.Println("Simple Admin's environment variables")
		color.Red.Println("Notice: Environment variables have priority over configuration files")
	}

	switch ServiceName {
	case "core":
		toolEnvInfo()
		authEnvInfo()
		crosEnvInfo()
		apiEnvInfo()
		rpcEnvInfo()
		logEnvInfo()
		databaseEnvInfo()
		i18nEnvInfo()
		captchaEnvInfo()
	case "fms":
		fmsEnvInfo()
	}

	return nil
}
