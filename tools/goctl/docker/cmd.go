package docker

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

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
	Cmd = &cobra.Command{
		Use:   "docker",
		Short: "Generate Dockerfile",
		RunE:  dockerCommand,
	}
)

func init() {
	dockerCmdFlags := Cmd.Flags()
	dockerCmdFlags.StringVar(&varExeName, "exe", "", flags.Get("docker.exe"))
	dockerCmdFlags.StringVar(&varStringGo, "go", "", flags.Get("docker.go"))
	dockerCmdFlags.StringVar(&varStringBase, "base", "scratch", flags.Get("docker.base"))
	dockerCmdFlags.IntVar(&varIntPort, "port", 0, flags.Get("docker.port"))
	dockerCmdFlags.StringVar(&varStringHome, "home", "", flags.Get("docker.home"))
	dockerCmdFlags.StringVar(&varStringRemote, "remote", "", flags.Get("docker.remote"))
	dockerCmdFlags.StringVar(&varStringBranch, "branch", "", flags.Get("docker.branch"))
	dockerCmdFlags.StringVar(&varStringVersion, "version", "", flags.Get("docker.version"))
	dockerCmdFlags.StringVar(&varStringTZ, "tz", "Asia/Shanghai", flags.Get("docker.tz"))
}
