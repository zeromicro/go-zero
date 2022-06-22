package docker

import "github.com/spf13/cobra"

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
	Cmd.Flags().StringVar(&varExeName, "exe", "", "The executable name in the built image")
	Cmd.Flags().StringVar(&varStringGo, "go", "", "The file that contains main function")
	Cmd.Flags().StringVar(&varStringBase, "base", "scratch", "The base image to build the docker image, default scratch")
	Cmd.Flags().IntVar(&varIntPort, "port", 0, "The port to expose, default none")
	Cmd.Flags().StringVar(&varStringHome, "home", "", "The goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	Cmd.Flags().StringVar(&varStringRemote, "remote", "", "The remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\nThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	Cmd.Flags().StringVar(&varStringBranch, "branch", "", "The branch of the remote repo, it does work with --remote")
	Cmd.Flags().StringVar(&varStringVersion, "version", "", "The goctl builder golang image version")
	Cmd.Flags().StringVar(&varStringTZ, "tz", "Asia/Shanghai", "The timezone of the container")
}
