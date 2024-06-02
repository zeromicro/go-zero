package docker

import "github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"

var (
	varExeName       string
	varStringGo      string
	varStringBase    string
	varIntPort       int
	varStringHome    string
	varStringRemote  string
	varStringBranch  string
	varStringVersion string
	varStringTZ      string

	// Cmd describes a docker command.
	Cmd = cobrax.NewCommand("docker", cobrax.WithRunE(dockerCommand))
)

func init() {
	dockerCmdFlags := Cmd.Flags()
	dockerCmdFlags.StringVar(&varExeName, "exe")
	dockerCmdFlags.StringVar(&varStringGo, "go")
	dockerCmdFlags.StringVarWithDefaultValue(&varStringBase, "base", "scratch")
	dockerCmdFlags.IntVar(&varIntPort, "port")
	dockerCmdFlags.StringVar(&varStringHome, "home")
	dockerCmdFlags.StringVar(&varStringRemote, "remote")
	dockerCmdFlags.StringVar(&varStringBranch, "branch")
	dockerCmdFlags.StringVar(&varStringVersion, "version")
	dockerCmdFlags.StringVarWithDefaultValue(&varStringTZ, "tz", "Asia/Shanghai")
}
