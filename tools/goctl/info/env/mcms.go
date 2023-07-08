package env

import (
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func mcmsEmailEnvInfo() string {
	color.Green.Println("MCMS")
	color.Green.Println("EMAIL")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"EMAIL_AUTH_TYPE", "电子邮箱的认证类型，支持plain和CRAMMD5"},
			{"EMAIL_ADDR", "电子邮箱地址"},
			{"EMAIL_PASSWORD", "电子邮箱密码"},
			{"EMAIL_HOST_NAME", "电子邮箱的服务器地址"},
			{"EMAIL_PORT", "电子邮箱的服务器端口"},
			{"EMAIL_IDENTIFY", "电子邮箱的身份信息，用于CRAMMD5"},
			{"EMAIL_SECRET", "电子邮箱的密钥信息，用于CRAMMD5"},
			{"EMAIL_TLS", "是否启用TLS"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"EMAIL_AUTH_TYPE", "Authentication type for the email, supports plain and CRAMMD5"},
			{"EMAIL_ADDR", "Email address"},
			{"EMAIL_PASSWORD", "Email password"},
			{"EMAIL_HOST_NAME", "Server address for the email"},
			{"EMAIL_PORT", "Server port for the email"},
			{"EMAIL_IDENTIFY", "Identity information for the email, used for CRAMMD5"},
			{"EMAIL_SECRET", "Secret information for the email, used for CRAMMD5"},
			{"EMAIL_TLS", "Whether to enable TLS"},
		})
	}
	return envInfo.Render()
}

func mcmsSmsEnvInfo() string {
	color.Green.Println("SMS")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"SMS_SECRET_ID", "短信服务密钥ID"},
			{"SMS_SECRET_KEY", "短信服务密钥"},
			{"SMS_PROVIDER", "短信服务提供商"},
			{"SMS_REGION", "短信服务提供区域"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"SMS_SECRET_ID", "Secret ID for the SMS service"},
			{"SMS_SECRET_KEY", "Secret key for the SMS service"},
			{"SMS_PROVIDER", "Provider for the SMS service"},
			{"SMS_REGION", "Region for the SMS service"},
		})
	}
	return envInfo.Render()
}
