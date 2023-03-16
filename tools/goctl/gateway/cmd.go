package gateway

import (
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/cli"
)

var (
	varStringHome   string
	varStringRemote string
	varStringBranch string

	Cmd = &cobra.Command{
		Use:   "gateway",
		Short: "gateway is a tool to generate gateway code",
	}

	protoCmd = &cobra.Command{
		Use:   "proto",
		Short: "generate gateway code from proto file",
		RunE:  runFromProto,
	}

	protosetCmd = &cobra.Command{
		Use:   "protoset",
		Short: "generate gateway code from protoset file",
		RunE:  runFromProtoSet,
	}

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "generate gateway code from grpc server",
		RunE:  runFromGRPCServer,
	}
)

func init() {
	Cmd.PersistentFlags().StringVar(&cli.VarStringHome, "home", "", "The goctl home"+
		" path of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority")
	Cmd.PersistentFlags().StringVar(&cli.VarStringRemote, "remote", "", "The remote "+
		"git repo of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	Cmd.PersistentFlags().StringVar(&cli.VarStringBranch, "branch", "", "The branch"+
		" of the remote repo, it does work with --remote")

	Cmd.AddCommand(protoCmd, protosetCmd, serverCmd)
}
