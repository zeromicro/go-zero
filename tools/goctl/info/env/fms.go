package env

import (
	"os"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

func fmsEnvInfo() string {
	color.Green.Println("FMS")
	envInfo = table.NewWriter()
	envInfo.SetOutputMirror(os.Stdout)
	if lang {
		envInfo.AppendHeader(table.Row{"环境变量名称", "环境变量介绍"})
		envInfo.AppendRows([]table.Row{
			{"MAX_IMAGE_SIZE", "图片文件最大大小"},
			{"MAX_VIDEO_SIZE", "视频文件最大大小"},
			{"MAX_AUDIO_SIZE", "音频文件最大大小"},
			{"MAX_OTHER_SIZE", "其他类型文件最大大小"},
			{"PRIVATE_PATH", "私人文件本地存储地址"},
			{"PUBLIC_PATH", "公开文件本地存储地址"},
			{"SERVER_URL", "服务器域名或IP地址，如 http://localhost:81"},
		})
	} else {
		envInfo.AppendHeader(table.Row{"Key", "Introduction"})
		envInfo.AppendRows([]table.Row{
			{"MAX_IMAGE_SIZE", "Maximum size of image files"},
			{"MAX_VIDEO_SIZE", "Maximum size of video files"},
			{"MAX_AUDIO_SIZE", "Maximum size of audio files"},
			{"MAX_OTHER_SIZE", "Maximum size of files of other types"},
			{"PRIVATE_PATH", "Local storage address for private files"},
			{"PUBLIC_PATH", "Local storage address for public files"},
			{"SERVER_URL", "Domain name or IP address of the server, such as http://localhost:81"},
		})
	}
	return envInfo.Render()
}
