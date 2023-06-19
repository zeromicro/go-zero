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
	"os"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/zeromicro/go-zero/tools/goctl/util/env"
)

var envInfo table.Writer
var lang = env.IsChinaEnv()

// toolEnvInfo show the tools env variables usage by goctls
func toolEnvInfo() string {
	color.Green.Println("TOOLS")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"SIMPLE_ADMIN_TOOLS_LANG", "控制台中goctls的帮助信息语言类型，支持zh和en，默认为en"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"SIMPLE_ADMIN_TOOLS_LANG", "The goctls' help message language type in console, support zh and en, default is en"},
		})
	}
	return envInfo.Render()
}

// serviceEnvInfo show the api env variables usage by goctls
func apiEnvInfo() string {
	color.Green.Println("API")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"API_HOST", "API服务的主机地址"},
			{"API_PORT", "API 服务端口"},
			{"API_TIMEOUT", "API 服务超时时间"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"API_HOST", "The API service's host address"},
			{"API_PORT", "The API service's port"},
			{"API_TIMEOUT", "The API service's timeout"},
		})
	}
	return envInfo.Render()
}

// rpcEnvInfo show the rpc env variables usage by goctls
func rpcEnvInfo() string {
	color.Green.Println("RPC")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"RPC_LISTEN_ON", "RPC服务的主机和端口地址，如localhost:8080"},
			{"RPC_TIMEOUT", "RPC 服务的超时设置"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"RPC_LISTEN_ON", "The RPC service's host and port address, such as localhost:8080"},
			{"RPC_TIMEOUT", "The RPC service's timeout setting"},
		})
	}
	return envInfo.Render()
}

// logEnvInfo show the log env variables usage by goctls
func logEnvInfo() string {
	color.Green.Println("LOG")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"LOG_SERVICE_NAME", "日志中的服务名称"},
			{"LOG_MODE", "日志模式, 如 file, console 和 volume"},
			{"LOG_ENCODING", "日志编码如 json、plain"},
			{"LOG_PATH", "日志存放路径，当 mode 为 file 时启用"},
			{"LOG_LEVEL", "日志级别如 info, error"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"LOG_SERVICE_NAME", "The service's name in log"},
			{"LOG_MODE", "The log mode such as file, console and volume"},
			{"LOG_ENCODING", "The log encoding such as json and plain"},
			{"LOG_PATH", "The log storage path, use in file mode"},
			{"LOG_LEVEL", "The log level such as info and error"},
		})
	}
	return envInfo.Render()
}

// databaseEnvInfo show the database env variables usage by goctls
func databaseEnvInfo() string {
	color.Green.Println("DATABASE (Ent)")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"DATABASE_HOST", "数据库的主机地址"},
			{"DATABASE_PORT", "数据库端口"},
			{"DATABASE_USERNAME", "数据库的用户名"},
			{"DATABASE_PASSWORD", "数据库密码"},
			{"DATABASE_DBNAME", "数据库名称"},
			{"DATABASE_SSL_MODE", "数据库的ssl模式"},
			{"DATABASE_TYPE", "数据库类型，支持mysql、postgres和sqlite3"},
			{"DATABASE_MAX_OPEN_CONN", "数据库的最大打开连接数"},
			{"DATABASE_CACHE_TIME", "数据库的缓存时间"},
			{"DATABASE_DBPATH", "sqlite3 的数据库存放路径"},
			{"DATABASE_MYSQL_CONFIG", "数据库对 mysql 的额外配置"},
			{"DATABASE_PG_CONFIG", "数据库对 postgresql 的额外配置"},
			{"DATABASE_SQLITE_CONFIG", "数据库对 sqlite3 的额外配置"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"DATABASE_HOST", "The database's host address"},
			{"DATABASE_PORT", "The database's port"},
			{"DATABASE_USERNAME", "The database's username"},
			{"DATABASE_PASSWORD", "The database's password"},
			{"DATABASE_DBNAME", "The database's name"},
			{"DATABASE_SSL_MODE", "The database's ssl mode"},
			{"DATABASE_TYPE", "The database's type, support mysql, postgres and sqlite3"},
			{"DATABASE_MAX_OPEN_CONN", "The database's max opened connections"},
			{"DATABASE_CACHE_TIME", "The database's cache time"},
			{"DATABASE_DBPATH", "The database's storage path for sqlite3"},
			{"DATABASE_MYSQL_CONFIG", "The database's extra config for mysql"},
			{"DATABASE_PG_CONFIG", "The database's extra config for postgresql"},
			{"DATABASE_SQLITE_CONFIG", "The database's extra config for sqlite3"},
		})
	}
	return envInfo.Render()
}

// captchaEnvInfo show the captcha env variables usage by goctls
func captchaEnvInfo() string {
	color.Green.Println("CAPTCHA")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"CAPTCHA_KEY_LONG", "验证码长度"},
			{"CAPTCHA_IMG_WIDTH", "验证码宽度"},
			{"CAPTCHA_IMG_HEIGHT", "验证码高度"},
			{"CAPTCHA_DRIVER", "验证码类型如 math, string"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"CAPTCHA_KEY_LONG", "The captcha key length"},
			{"CAPTCHA_IMG_WIDTH", "The width of captcha"},
			{"CAPTCHA_IMG_HEIGHT", "The height of captcha"},
			{"CAPTCHA_DRIVER", "The driver of captcha such as math, string and digit"},
		})
	}
	return envInfo.Render()
}

// authEnvInfo show the auth env variables usage by goctls
func authEnvInfo() string {
	color.Green.Println("JWT")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"AUTH_SECRET", "JWT加密密钥"},
			{"AUTH_EXPIRE", "JWT过期时间"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"AUTH_SECRET", "JWT encrypted key"},
			{"AUTH_EXPIRE", "JWT expired time"},
		})
	}
	return envInfo.Render()
}

// i18nEnvInfo show the i18n env variables usage by goctls
func i18nEnvInfo() string {
	color.Green.Println("I18n")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"I18N_DIR", "i18n 外部文件目录，即包含 locale 文件夹的目录"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"I18N_DIR", "The i18n external file directory, that is the directory containing the locale folder"},
		})
	}
	return envInfo.Render()
}

// crosEnvInfo show the cros env variables usage by goctls
func crosEnvInfo() string {
	color.Green.Println("CROS")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"CROS_ADDRESS", "跨域允许的域名或 ip, 如 http://qq.com"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"CROS_ADDRESS", "Domain name or ip allowed across domains, such as http://qq.com"},
		})
	}
	return envInfo.Render()
}
