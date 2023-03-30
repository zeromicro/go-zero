package api

import (
	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/api/apigen"
	"github.com/zeromicro/go-zero/tools/goctl/api/dartgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/docgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/format"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/api/javagen"
	"github.com/zeromicro/go-zero/tools/goctl/api/ktgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/new"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen"
	"github.com/zeromicro/go-zero/tools/goctl/api/validate"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

var (
	// Cmd describes an api command.
	Cmd       = cobrax.NewCommand("api", cobrax.WithRunE(apigen.CreateApiTemplate))
	dartCmd   = cobrax.NewCommand("dart", cobrax.WithRunE(dartgen.DartCommand))
	docCmd    = cobrax.NewCommand("doc", cobrax.WithRunE(docgen.DocCommand))
	formatCmd = cobrax.NewCommand("format", cobrax.WithRunE(format.GoFormatApi))
	goCmd     = cobrax.NewCommand("go", cobrax.WithRunE(gogen.GoCommand))
	newCmd    = cobrax.NewCommand("new", cobrax.WithRunE(new.CreateServiceCommand),
		cobrax.WithArgs(cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)))
	validateCmd = cobrax.NewCommand("validate", cobrax.WithRunE(validate.GoValidateApi))
	javaCmd     = cobrax.NewCommand("java", cobrax.WithRunE(javagen.JavaCommand), cobrax.WithHidden())
	ktCmd       = cobrax.NewCommand("kt", cobrax.WithRunE(ktgen.KtCommand))
	pluginCmd   = cobrax.NewCommand("plugin", cobrax.WithRunE(plugin.PluginCommand))
	tsCmd       = cobrax.NewCommand("ts", cobrax.WithRunE(tsgen.TsCommand))

	protoCmd = cobrax.NewCommand("proto", cobrax.WithRunE(gogen.GenCRUDLogicByProto))

	entCmd = cobrax.NewCommand("ent", cobrax.WithRunE(gogen.GenCRUDLogicByEnt))
)

func init() {
	var (
		apiCmdFlags      = Cmd.Flags()
		dartCmdFlags     = dartCmd.Flags()
		docCmdFlags      = docCmd.Flags()
		formatCmdFlags   = formatCmd.Flags()
		goCmdFlags       = goCmd.Flags()
		javaCmdFlags     = javaCmd.Flags()
		ktCmdFlags       = ktCmd.Flags()
		newCmdFlags      = newCmd.Flags()
		pluginCmdFlags   = pluginCmd.Flags()
		tsCmdFlags       = tsCmd.Flags()
		validateCmdFlags = validateCmd.Flags()
		protoCmdFlags    = protoCmd.Flags()
		entCmdFlags      = entCmd.Flags()
	)

	apiCmdFlags.StringVar(&apigen.VarStringOutput, "o")
	apiCmdFlags.StringVar(&apigen.VarStringHome, "home")
	apiCmdFlags.StringVar(&apigen.VarStringRemote, "remote")
	apiCmdFlags.StringVar(&apigen.VarStringBranch, "branch")

	dartCmdFlags.StringVar(&dartgen.VarStringDir, "dir")
	dartCmdFlags.StringVar(&dartgen.VarStringAPI, "api")
	dartCmdFlags.BoolVar(&dartgen.VarStringLegacy, "legacy")
	dartCmdFlags.StringVar(&dartgen.VarStringHostname, "hostname")
	dartCmdFlags.StringVar(&dartgen.VarStringScheme, "scheme")

	docCmdFlags.StringVar(&docgen.VarStringDir, "dir")
	docCmdFlags.StringVar(&docgen.VarStringOutput, "o")

	formatCmdFlags.StringVar(&format.VarStringDir, "dir")
	formatCmdFlags.BoolVar(&format.VarBoolIgnore, "iu")
	formatCmdFlags.BoolVar(&format.VarBoolUseStdin, "stdin")
	formatCmdFlags.BoolVar(&format.VarBoolSkipCheckDeclare, "declare")

	goCmdFlags.StringVar(&gogen.VarStringDir, "dir")
	goCmdFlags.StringVar(&gogen.VarStringAPI, "api")
	goCmdFlags.StringVar(&gogen.VarStringHome, "home")
	goCmdFlags.StringVar(&gogen.VarStringRemote, "remote")
	goCmdFlags.StringVar(&gogen.VarStringBranch, "branch")
	goCmdFlags.StringVarWithDefaultValue(&gogen.VarStringStyle, "style", config.DefaultFormat)
	goCmdFlags.BoolVar(&gogen.VarBoolErrorTranslate, "trans_err")
	goCmdFlags.BoolVar(&gogen.VarBoolUseCasbin, "casbin")
	goCmdFlags.BoolVar(&gogen.VarBoolUseI18n, "i18n")

	javaCmdFlags.StringVar(&javagen.VarStringDir, "dir")
	javaCmdFlags.StringVar(&javagen.VarStringAPI, "api")

	ktCmdFlags.StringVar(&ktgen.VarStringDir, "dir")
	ktCmdFlags.StringVar(&ktgen.VarStringAPI, "api")
	ktCmdFlags.StringVar(&ktgen.VarStringPKG, "pkg")

	newCmdFlags.StringVar(&new.VarStringHome, "home")
	newCmdFlags.StringVar(&new.VarStringRemote, "remote")
	newCmdFlags.StringVar(&new.VarStringBranch, "branch")
	newCmdFlags.StringVarWithDefaultValue(&new.VarStringStyle, "style", config.DefaultFormat)
	newCmdFlags.BoolVar(&new.VarBoolUseCasbin, "casbin")
	newCmdFlags.BoolVar(&new.VarBoolUseI18n, "i18n")
	newCmdFlags.StringVar(&new.VarStringGoZeroVersion, "go_zero_version")
	newCmdFlags.StringVar(&new.VarStringToolVersion, "tool_version")
	newCmdFlags.StringVar(&new.VarModuleName, "module_name")
	newCmdFlags.BoolVar(&new.VarBoolErrorTranslate, "trans_err")
	newCmdFlags.IntVarWithDefaultValue(&new.VarIntServicePort, "port", 9100)
	newCmdFlags.BoolVar(&new.VarBoolGitlab, "gitlab")
	newCmdFlags.BoolVar(&new.VarBoolEnt, "ent")

	pluginCmdFlags.StringVarP(&plugin.VarStringPlugin, "plugin", "p")
	pluginCmdFlags.StringVar(&plugin.VarStringDir, "dir")
	pluginCmdFlags.StringVar(&plugin.VarStringAPI, "api")
	pluginCmdFlags.StringVar(&plugin.VarStringStyle, "style")

	tsCmdFlags.StringVar(&tsgen.VarStringDir, "dir")
	tsCmdFlags.StringVar(&tsgen.VarStringAPI, "api")
	tsCmdFlags.StringVar(&tsgen.VarStringCaller, "caller")
	tsCmdFlags.BoolVar(&tsgen.VarBoolUnWrap, "unwrap")

	validateCmdFlags.StringVar(&validate.VarStringAPI, "api")

	protoCmdFlags.StringVar(&gogen.VarStringProto, "proto")
	protoCmdFlags.StringVar(&gogen.VarStringOutput, "o")
	protoCmdFlags.StringVar(&gogen.VarStringAPIServiceName, "api_service_name")
	protoCmdFlags.StringVar(&gogen.VarStringRPCServiceName, "rpc_service_name")
	protoCmdFlags.StringVarWithDefaultValue(&gogen.VarStringStyle, "style", config.DefaultFormat)
	protoCmdFlags.StringVar(&gogen.VarStringModelName, "model")
	protoCmdFlags.IntVarWithDefaultValue(&gogen.VarIntSearchKeyNum, "search_key_num", 3)
	protoCmdFlags.StringVar(&gogen.VarStringRpcName, "rpc_name")
	protoCmdFlags.StringVar(&gogen.VarStringGrpcPbPackage, "grpc_package")
	protoCmdFlags.BoolVar(&gogen.VarBoolMultiple, "multiple")
	protoCmdFlags.StringVarWithDefaultValue(&gogen.VarStringJSONStyle, "json_style", "goZero")
	protoCmdFlags.BoolVar(&gogen.VarBoolOverwrite, "overwrite")

	entCmdFlags.StringVar(&gogen.VarStringSchema, "schema")
	entCmdFlags.StringVar(&gogen.VarStringOutput, "o")
	entCmdFlags.StringVar(&gogen.VarStringAPIServiceName, "api_service_name")
	entCmdFlags.StringVarWithDefaultValue(&gogen.VarStringStyle, "style", config.DefaultFormat)
	entCmdFlags.StringVar(&gogen.VarStringModelName, "model")
	entCmdFlags.IntVarWithDefaultValue(&gogen.VarIntSearchKeyNum, "search_key_num", 3)
	entCmdFlags.StringVar(&gogen.VarStringGroupName, "group")
	entCmdFlags.BoolVar(&gogen.VarBoolOverwrite, "overwrite")
	entCmdFlags.StringVarWithDefaultValue(&gogen.VarStringJSONStyle, "json_style", "goZero")

	// Add sub-commands
	Cmd.AddCommand(dartCmd)
	Cmd.AddCommand(docCmd)
	Cmd.AddCommand(formatCmd)
	Cmd.AddCommand(goCmd)
	Cmd.AddCommand(javaCmd)
	Cmd.AddCommand(ktCmd)
	Cmd.AddCommand(newCmd)
	Cmd.AddCommand(pluginCmd)
	Cmd.AddCommand(tsCmd)
	Cmd.AddCommand(validateCmd)
	Cmd.AddCommand(protoCmd)
	Cmd.AddCommand(entCmd)
}
