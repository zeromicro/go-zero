package rpc

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/cli"
)

var (
	// Cmd describes a rpc command.
	Cmd = &cobra.Command{
		Use:   "rpc",
		Short: "Generate rpc code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.RPCTemplate(true)
		},
	}

	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Generate rpc demo service",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:  cli.RPCNew,
	}

	templateCmd = &cobra.Command{
		Use:   "template",
		Short: "Generate proto template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.RPCTemplate(false)
		},
	}

	protocCmd = &cobra.Command{
		Use:     "protoc",
		Short:   "Generate grpc code",
		Example: "goctl rpc protoc xx.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:    cli.ZRPC,
	}
)

func init() {
	Cmd.Flags().StringVar(&cli.VarStringOutput, "o", "", "Output a sample proto file")
	Cmd.Flags().StringVar(&cli.VarStringHome, "home", "", "The goctl home path of "+
		"the template, --home and --remote cannot be set at the same time, if they are, --remote has"+
		" higher priority")
	Cmd.Flags().StringVar(&cli.VarStringRemote, "remote", "", "The remote git repo"+
		" of the template, --home and --remote cannot be set at the same time, if they are, --remote"+
		" has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	Cmd.Flags().StringVar(&cli.VarStringBranch, "branch", "", "The branch of the "+
		"remote repo, it does work with --remote")

	newCmd.Flags().StringSliceVar(&cli.VarStringSliceGoOpt, "go_opt", nil, "")
	newCmd.Flags().StringSliceVar(&cli.VarStringSliceGoGRPCOpt, "go-grpc_opt", nil, "")
	newCmd.Flags().StringVar(&cli.VarStringStyle, "style", "gozero", "The file "+
		"naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	newCmd.Flags().BoolVar(&cli.VarBoolIdea, "idea", false, "Whether the command "+
		"execution environment is from idea plugin.")
	newCmd.Flags().StringVar(&cli.VarStringHome, "home", "", "The goctl home path "+
		"of the template, --home and --remote cannot be set at the same time, if they are, --remote "+
		"has higher priority")
	newCmd.Flags().StringVar(&cli.VarStringRemote, "remote", "", "The remote git "+
		"repo of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	newCmd.Flags().StringVar(&cli.VarStringBranch, "branch", "",
		"The branch of the remote repo, it does work with --remote")
	newCmd.Flags().BoolVarP(&cli.VarBoolVerbose, "verbose", "v", false, "Enable log output")
	newCmd.Flags().MarkHidden("go_opt")
	newCmd.Flags().MarkHidden("go-grpc_opt")

	protocCmd.Flags().BoolVarP(&cli.VarBoolMultiple, "multiple", "m", false,
		"Generated in multiple rpc service mode")
	protocCmd.Flags().StringSliceVar(&cli.VarStringSliceGoOut, "go_out", nil, "")
	protocCmd.Flags().StringSliceVar(&cli.VarStringSliceGoGRPCOut, "go-grpc_out", nil, "")
	protocCmd.Flags().StringSliceVar(&cli.VarStringSliceGoOpt, "go_opt", nil, "")
	protocCmd.Flags().StringSliceVar(&cli.VarStringSliceGoGRPCOpt, "go-grpc_opt", nil, "")
	protocCmd.Flags().StringSliceVar(&cli.VarStringSlicePlugin, "plugin", nil, "")
	protocCmd.Flags().StringSliceVarP(&cli.VarStringSliceProtoPath, "proto_path", "I", nil, "")
	protocCmd.Flags().StringVar(&cli.VarStringZRPCOut, "zrpc_out", "", "The zrpc output directory")
	protocCmd.Flags().StringVar(&cli.VarStringStyle, "style", "gozero", "The file "+
		"naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")
	protocCmd.Flags().StringVar(&cli.VarStringHome, "home", "", "The goctl home "+
		"path of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority")
	protocCmd.Flags().StringVar(&cli.VarStringRemote, "remote", "", "The remote "+
		"git repo of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	protocCmd.Flags().StringVar(&cli.VarStringBranch, "branch", "",
		"The branch of the remote repo, it does work with --remote")
	protocCmd.Flags().BoolVarP(&cli.VarBoolVerbose, "verbose", "v", false, "Enable log output")
	protocCmd.Flags().MarkHidden("go_out")
	protocCmd.Flags().MarkHidden("go-grpc_out")
	protocCmd.Flags().MarkHidden("go_opt")
	protocCmd.Flags().MarkHidden("go-grpc_opt")
	protocCmd.Flags().MarkHidden("plugin")
	protocCmd.Flags().MarkHidden("proto_path")

	templateCmd.Flags().StringVar(&cli.VarStringOutput, "o", "", "Output a sample proto file")
	templateCmd.Flags().StringVar(&cli.VarStringHome, "home", "", "The goctl home"+
		" path of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority")
	templateCmd.Flags().StringVar(&cli.VarStringRemote, "remote", "", "The remote "+
		"git repo of the template, --home and --remote cannot be set at the same time, if they are, "+
		"--remote has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	templateCmd.Flags().StringVar(&cli.VarStringBranch, "branch", "", "The branch"+
		" of the remote repo, it does work with --remote")

	Cmd.AddCommand(newCmd)
	Cmd.AddCommand(protocCmd)
	Cmd.AddCommand(templateCmd)
}
