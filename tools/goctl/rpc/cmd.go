package rpc

import (
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/cli"
)

var (
	// Cmd describes a rpc command.
	Cmd = cobrax.NewCommand("rpc", cobrax.WithRunE(func(command *cobra.Command, strings []string) error {
		return cli.RPCTemplate(true)
	}))
	templateCmd = cobrax.NewCommand("template", cobrax.WithRunE(func(command *cobra.Command, strings []string) error {
		return cli.RPCTemplate(false)
	}))

	newCmd    = cobrax.NewCommand("new", cobrax.WithRunE(cli.RPCNew), cobrax.WithArgs(cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)))
	protocCmd = cobrax.NewCommand("protoc", cobrax.WithRunE(cli.ZRPC), cobrax.WithArgs(cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)))

	entCmd = cobrax.NewCommand("ent", cobrax.WithRunE(cli.EntCRUDLogic))
)

func init() {
	var (
		rpcCmdFlags      = Cmd.Flags()
		newCmdFlags      = newCmd.Flags()
		protocCmdFlags   = protocCmd.Flags()
		templateCmdFlags = templateCmd.Flags()
		entCmdFlags      = entCmd.Flags()
	)

	rpcCmdFlags.StringVar(&cli.VarStringOutput, "o")
	rpcCmdFlags.StringVar(&cli.VarStringHome, "home")
	rpcCmdFlags.StringVar(&cli.VarStringRemote, "remote")
	rpcCmdFlags.StringVar(&cli.VarStringBranch, "branch")

	newCmdFlags.StringSliceVar(&cli.VarStringSliceGoOpt, "go_opt")
	newCmdFlags.StringSliceVar(&cli.VarStringSliceGoGRPCOpt, "go-grpc_opt")
	newCmdFlags.StringVarWithDefaultValue(&cli.VarStringStyle, "style", config.DefaultFormat)
	newCmdFlags.BoolVar(&cli.VarBoolIdea, "idea")
	newCmdFlags.StringVar(&cli.VarStringHome, "home")
	newCmdFlags.StringVar(&cli.VarStringRemote, "remote")
	newCmdFlags.StringVar(&cli.VarStringBranch, "branch")
	newCmdFlags.BoolVarP(&cli.VarBoolVerbose, "verbose", "v")
	newCmdFlags.MarkHidden("go_opt")
	newCmdFlags.MarkHidden("go-grpc_opt")
	newCmdFlags.BoolVar(&cli.VarBoolEnt, "ent")
	newCmdFlags.StringVar(&cli.VarStringModuleName, "module_name")
	newCmdFlags.StringVar(&cli.VarStringGoZeroVersion, "go_zero_version")
	newCmdFlags.StringVar(&cli.VarStringToolVersion, "tool_version")
	newCmdFlags.IntVarWithDefaultValue(&cli.VarIntServicePort, "port", 9110)
	newCmdFlags.BoolVar(&cli.VarBoolGitlab, "gitlab")
	newCmdFlags.BoolVar(&cli.VarBoolDesc, "desc")

	protocCmdFlags.BoolVarP(&cli.VarBoolMultiple, "multiple", "m")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoOut, "go_out")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoGRPCOut, "go-grpc_out")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoOpt, "go_opt")
	protocCmdFlags.StringSliceVar(&cli.VarStringSliceGoGRPCOpt, "go-grpc_opt")
	protocCmdFlags.StringSliceVar(&cli.VarStringSlicePlugin, "plugin")
	protocCmdFlags.StringSliceVarP(&cli.VarStringSliceProtoPath, "proto_path", "I")
	protocCmdFlags.StringVar(&cli.VarStringZRPCOut, "zrpc_out")
	protocCmdFlags.StringVar(&cli.VarStringHome, "home")
	protocCmdFlags.StringVar(&cli.VarStringRemote, "remote")
	protocCmdFlags.StringVar(&cli.VarStringBranch, "branch")
	protocCmdFlags.BoolVarP(&cli.VarBoolVerbose, "verbose", "v")
	protocCmdFlags.MarkHidden("go_out")
	protocCmdFlags.MarkHidden("go-grpc_out")
	protocCmdFlags.MarkHidden("go_opt")
	protocCmdFlags.MarkHidden("go-grpc_opt")
	protocCmdFlags.MarkHidden("plugin")
	protocCmdFlags.MarkHidden("proto_path")

	templateCmdFlags.StringVar(&cli.VarStringOutput, "o")
	templateCmdFlags.StringVar(&cli.VarStringHome, "home")
	templateCmdFlags.StringVar(&cli.VarStringRemote, "remote")
	templateCmdFlags.StringVar(&cli.VarStringBranch, "branch")

	entCmdFlags.StringVar(&cli.VarStringSchema, "schema")
	entCmdFlags.StringVar(&cli.VarStringOutput, "o")
	entCmdFlags.StringVar(&cli.VarStringServiceName, "service_name")
	entCmdFlags.StringVar(&cli.VarStringProjectName, "project_name")
	entCmdFlags.BoolVar(&cli.VarBoolMultiple, "multiple")
	entCmdFlags.StringVarWithDefaultValue(&cli.VarStringStyle, "style", config.DefaultFormat)
	entCmdFlags.StringVar(&cli.VarStringModelName, "model")
	entCmdFlags.IntVarWithDefaultValue(&cli.VarIntSearchKeyNum, "search_key_num", 3)
	entCmdFlags.StringVar(&cli.VarStringGroupName, "group")
	entCmdFlags.StringVar(&cli.VarStringProtoPath, "proto_out")
	entCmdFlags.BoolVar(&cli.VarBoolOverwrite, "overwrite")

	Cmd.AddCommand(newCmd)
	Cmd.AddCommand(protocCmd)
	Cmd.AddCommand(templateCmd)
	Cmd.AddCommand(entCmd)
}
