package docker

import "github.com/spf13/cobra"

var (
	varStringGo      string
	varStringBase    string
	varIntPort       int
	varStringHome    string
	varStringRemote  string
	varStringBranch  string
	varStringVersion string
	varStringTZ      string
	Cmd              = &cobra.Command{
		Use:   "docker",
		Short: "generate Dockerfile",
		RunE:  dockerCommand,
	}
)

func init() {
	Cmd.Flags().StringVar(&varStringGo, "go", "", "the file that contains main function")
	Cmd.Flags().StringVar(&varStringBase, "base", "scratch", "the base image to build the docker image, default scratch")
	Cmd.Flags().IntVar(&varIntPort, "port", 0, "the port to expose, default none")
	Cmd.Flags().StringVar(&varStringHome, "home", "", "the goctl home path of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	Cmd.Flags().StringVar(&varStringRemote, "remote", "", "the remote git repo of the template, --home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	Cmd.Flags().StringVar(&varStringBranch, "branch", "", "the branch of the remote repo, it does work with --remote")
	Cmd.Flags().StringVar(&varStringVersion, "version", "", "the goctl builder golang image version")
	Cmd.Flags().StringVar(&varStringTZ, "tz", "Asia/Shanghai", "the timezone of the container")
}
