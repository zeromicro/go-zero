package rpc

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/cli"
)

var (
	// Cmd describes a rpc command.
	Cmd = &cobra.Command{
		Use:   "rpc",
		Short: flags.Get("rpc.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.RPCTemplate(true)
		},
	}

	newCmd = &cobra.Command{
		Use:   "new",
		Short: flags.Get("rpc.new.short"),
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:  cli.RPCNew,
	}

	templateCmd = &cobra.Command{
		Use:   "template",
		Short: flags.Get("rpc.template.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.RPCTemplate(false)
		},
	}

	protocCmd = &cobra.Command{
		Use:     "protoc",
		Short:   flags.Get("rpc.protoc.short"),
		Example: flags.Get("rpc.protoc.example"),
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE:    cli.ZRPC,
	}
)

func init() {
	var (
		rpcCmdFlags      = Cmd.Flags()
		newCmdFlags      = newCmd.Flags()
		protocCmdFlags   = protocCmd.Flags()
		templateCmdFlags = templateCmd.Flags()
	)

	rpcCmdFlags.StringVar(&cli.VarStringOutput, "o", "", flags.Get("rpc.o"))
	rpcCmdFlags.StringVar(&cli.VarStringHome, "home", "", flags.Get("rpc.home"))
	rpcCmdFlags.StringVar(&cli.VarStringRemote, "remote", "", flags.Get("rpc.remote"))
	rpcCmdFlags.StringVar(&cli.VarStringBranch, "branch", "", flags.Get("rpc.branch"))

	newCmdFlags.StringSliceVar(&cli.VarStringSliceGoOpt, "go_opt", nil, "")
	newCmdFlags.StringSliceVar(&cli.VarStringSliceGoGRPCOpt, "go-grpc_opt", nil, "")
	newCmdFlags.StringVar(&cli.VarStringStyle, "style", config.DefaultFormat, flags.Get("rpc.new.style"))
	newCmdFlags.BoolVar(&cli.VarBoolIdea, "idea", false, flags.Get("rpc.new.idea"))
	newCmdFlags.StringVar(&cli.VarStringHome, "home", "", flags.Get("rpc.new..home"))
	newCmdFlags.StringVar(&cli.VarStringRemote, "remote", "", flags.Get("rpc.new.remote"))
	newCmdFlags.StringVar(&cli.VarStringBranch, "branch", "", flags.Get("rpc.new.branch"))
	newCmdFlags.BoolVarP(&cli.VarBoolVerbose, "verbose", "v", false, flags.Get("rpc.new.verbose"))
	newCmdFlags.MarkHidden("go_opt")
	newCmdFlags.MarkHidden("go-grpc_opt")

	protocCmdFlags.BoolVarP(&cli.VarBoolMultiple, "multiple", "m", false, flags.Get("rpc.protoc.multiple"))
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoOut, "go_out", nil, "")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoGRPCOut, "go-grpc_out", nil, "")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoOpt, "go_opt", nil, "")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoGRPCOpt, "go-grpc_opt", nil, "")
	protocCmdFlags.StringSliceVar(&cli.VarStringSlicePlugin, "plugin", nil, "")
	protocCmdFlags.StringSliceVarP(&cli.VarStringSliceProtoPath, "proto_path", "I", nil, "")
	protocCmdFlags.StringVar(&cli.VarStringZRPCOut, "zrpc_out", "", flags.Get("rpc.protoc.zrpc_out"))
	protocCmdFlags.StringVar(&cli.VarStringHome, "home", "", flags.Get("rpc.protoc.home"))
	protocCmdFlags.StringVar(&cli.VarStringRemote, "remote", "", flags.Get("rpc.protoc.remote"))
	protocCmdFlags.StringVar(&cli.VarStringBranch, "branch", "", flags.Get("rpc.protoc.branch"))
	protocCmdFlags.BoolVarP(&cli.VarBoolVerbose, "verbose", "v", false, flags.Get("rpc.protoc.verbose"))
	protocCmdFlags.MarkHidden("go_out")
	protocCmdFlags.MarkHidden("go-grpc_out")
	protocCmdFlags.MarkHidden("go_opt")
	protocCmdFlags.MarkHidden("go-grpc_opt")
	protocCmdFlags.MarkHidden("plugin")
	protocCmdFlags.MarkHidden("proto_path")

	templateCmdFlags.StringVar(&cli.VarStringOutput, "o", "", flags.Get("rpc.template.o"))
	templateCmdFlags.StringVar(&cli.VarStringHome, "home", "", flags.Get("rpc.template.home"))
	templateCmdFlags.StringVar(&cli.VarStringRemote, "remote", "", flags.Get("rpc.template.remote"))
	templateCmdFlags.StringVar(&cli.VarStringBranch, "branch", "", flags.Get("rpc.template.branch"))

	Cmd.AddCommand(newCmd, protocCmd, templateCmd)
}
