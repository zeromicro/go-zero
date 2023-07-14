package env

import (
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

// getServiceList returns the service list
func getServiceList() string {
	color.Green.Println("Service ")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"服务名称", "服务介绍"})
		envInfo.AppendRows([]table.Row{
			{"core", "核心服务"},
			{"fms", "文件服务"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Service name", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"core", "Core Service"},
			{"fms", "File Management Service"},
		})
	}
	return envInfo.Render()
}
