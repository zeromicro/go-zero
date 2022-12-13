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
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

var (
	// Cmd describes an api command.
	Cmd = &cobra.Command{
		Use:   "api",
		Short: "Generate api related files",
		RunE:  apigen.CreateApiTemplate,
	}

	dartCmd = &cobra.Command{
		Use:   "dart",
		Short: "Generate dart files for provided api in api file",
		RunE:  dartgen.DartCommand,
	}

	docCmd = &cobra.Command{
		Use:   "doc",
		Short: "Generate doc files",
		RunE:  docgen.DocCommand,
	}

	formatCmd = &cobra.Command{
		Use:   "format",
		Short: "Format api files",
		RunE:  format.GoFormatApi,
	}

	goCmd = &cobra.Command{
		Use:   "go",
		Short: "Generate go files for provided api in api file",
		RunE:  gogen.GoCommand,
	}

	newCmd = &cobra.Command{
		Use:     "new",
		Short:   "Fast create api service",
		Example: "goctl api new [options] service-name",
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			return new.CreateServiceCommand(args)
		},
	}

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Validate api file",
		RunE:  validate.GoValidateApi,
	}

	javaCmd = &cobra.Command{
		Use:    "java",
		Short:  "Generate java files for provided api in api file",
		Hidden: true,
		RunE:   javagen.JavaCommand,
	}

	ktCmd = &cobra.Command{
		Use:   "kt",
		Short: "Generate kotlin code for provided api file",
		RunE:  ktgen.KtCommand,
	}

	pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "Custom file generator",
		RunE:  plugin.PluginCommand,
	}

	tsCmd = &cobra.Command{
		Use:   "ts",
		Short: "Generate ts files for provided api in api file",
		RunE:  tsgen.TsCommand,
	}

	protoCmd = &cobra.Command{
		Use:   "proto",
		Short: "Generate CRUD template from proto file",
		RunE:  gogen.GenCRUDLogicByProto,
	}
)

func init() {
	Cmd.Flags().StringVar(&apigen.VarStringOutput, "o", "", "Output a sample api file")
	Cmd.Flags().StringVar(&apigen.VarStringHome, "home", "", "The goctl home path of the"+
		" template, --home and --remote cannot be set at the same time, if they are, --remote has "+
		"higher priority")
	Cmd.Flags().StringVar(&apigen.VarStringRemote, "remote", "", "The remote git repo of the"+
		" template, --home and --remote cannot be set at the same time, if they are, --remote has higher"+
		" priority\nThe git repo directory must be consistent with the"+
		" https://github.com/zeromicro/go-zero-template directory structure")
	Cmd.Flags().StringVar(&apigen.VarStringBranch, "branch", "", "The branch of the "+
		"remote repo, it does work with --remote")

	dartCmd.Flags().StringVar(&dartgen.VarStringDir, "dir", "", "The target dir")
	dartCmd.Flags().StringVar(&dartgen.VarStringAPI, "api", "", "The api file")
	dartCmd.Flags().BoolVar(&dartgen.VarStringLegacy, "legacy", false, "Legacy generator for flutter v1")
	dartCmd.Flags().StringVar(&dartgen.VarStringHostname, "hostname", "", "hostname of the server")

	docCmd.Flags().StringVar(&docgen.VarStringDir, "dir", "", "The target dir")
	docCmd.Flags().StringVar(&docgen.VarStringOutput, "o", "", "The output markdown directory")

	formatCmd.Flags().StringVar(&format.VarStringDir, "dir", "", "The format target dir")
	formatCmd.Flags().BoolVar(&format.VarBoolIgnore, "iu", false, "Ignore update")
	formatCmd.Flags().BoolVar(&format.VarBoolUseStdin, "stdin", false, "Use stdin to input api"+
		" doc content, press \"ctrl + d\" to send EOF")
	formatCmd.Flags().BoolVar(&format.VarBoolSkipCheckDeclare, "declare", false, "Use to skip check "+
		"api types already declare")

	goCmd.Flags().StringVar(&gogen.VarStringDir, "dir", "", "The target dir")
	goCmd.Flags().StringVar(&gogen.VarStringAPI, "api", "", "The api file")
	goCmd.Flags().StringVar(&gogen.VarStringHome, "home", "", "The goctl home path of "+
		"the template, --home and --remote cannot be set at the same time, if they are, --remote "+
		"has higher priority")
	goCmd.Flags().StringVar(&gogen.VarStringRemote, "remote", "", "The remote git repo "+
		"of the template, --home and --remote cannot be set at the same time, if they are, --remote"+
		" has higher priority\nThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	goCmd.Flags().StringVar(&gogen.VarStringBranch, "branch", "", "The branch of "+
		"the remote repo, it does work with --remote")
	goCmd.Flags().StringVar(&gogen.VarStringStyle, "style", "go_zero", "The file naming format,"+
		" see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
	goCmd.Flags().BoolVar(&gogen.VarBoolErrorTranslate, "transErr", false, "Whether to translate the error")
	goCmd.Flags().BoolVar(&gogen.VarBoolUseCasbin, "casbin", false, "Whether to use the Casbin")
	goCmd.Flags().BoolVar(&gogen.VarBoolUseI18n, "i18n", false, "Whether to use i18n")

	javaCmd.Flags().StringVar(&javagen.VarStringDir, "dir", "", "The target dir")
	javaCmd.Flags().StringVar(&javagen.VarStringAPI, "api", "", "The api file")

	ktCmd.Flags().StringVar(&ktgen.VarStringDir, "dir", "", "The target dir")
	ktCmd.Flags().StringVar(&ktgen.VarStringAPI, "api", "", "The api file")
	ktCmd.Flags().StringVar(&ktgen.VarStringPKG, "pkg", "", "Define package name for kotlin file")

	newCmd.Flags().StringVar(&new.VarStringHome, "home", "", "The goctl home path of "+
		"the template, --home and --remote cannot be set at the same time, if they are, --remote "+
		"has higher priority")
	newCmd.Flags().StringVar(&new.VarStringRemote, "remote", "", "The remote git repo "+
		"of the template, --home and --remote cannot be set at the same time, if they are, --remote"+
		" has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	newCmd.Flags().StringVar(&new.VarStringBranch, "branch", "", "The branch of "+
		"the remote repo, it does work with --remote")
	newCmd.Flags().StringVar(&new.VarStringStyle, "style", "go_zero", "The file naming format,"+
		" see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")
	newCmd.Flags().BoolVar(&new.VarBoolUseCasbin, "casbin", false, "Whether to use the Casbin")
	newCmd.Flags().BoolVar(&new.VarBoolUseI18n, "i18n", false, "Whether to use i18n")
	newCmd.Flags().StringVar(&new.VarStringGoZeroVersion, "goZeroVersion", "",
		"The go zero version used for migration. e.g. v1.4.2")
	newCmd.Flags().StringVar(&new.VarStringToolVersion, "toolVersion", "",
		"The simple admin tool version version used for migration. e.g. v0.0.9")
	newCmd.Flags().StringVar(&new.VarModuleName, "moduleName", "",
		"The module name in go.mod. e.g. github.com/suyuan32/simple-admin-core")
	newCmd.Flags().BoolVar(&new.VarBoolErrorTranslate, "transErr", false, "Whether to translate the error")
	newCmd.Flags().IntVar(&new.VarIntServicePort, "port", 9100, "The service port exposed")

	pluginCmd.Flags().StringVarP(&plugin.VarStringPlugin, "plugin", "p", "", "The plugin file")
	pluginCmd.Flags().StringVar(&plugin.VarStringDir, "dir", "", "The target dir")
	pluginCmd.Flags().StringVar(&plugin.VarStringAPI, "api", "", "The api file")
	pluginCmd.Flags().StringVar(&plugin.VarStringStyle, "style", "go_zero",
		"The file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")

	tsCmd.Flags().StringVar(&tsgen.VarStringDir, "dir", "", "The target dir")
	tsCmd.Flags().StringVar(&tsgen.VarStringAPI, "api", "", "The api file")
	tsCmd.Flags().StringVar(&tsgen.VarStringWebAPI, "webapi", "", "The web api file path")
	tsCmd.Flags().StringVar(&tsgen.VarStringCaller, "caller", "", "The web api caller")
	tsCmd.Flags().BoolVar(&tsgen.VarBoolUnWrap, "unwrap", false, "Unwrap the webapi caller for import")

	validateCmd.Flags().StringVar(&validate.VarStringAPI, "api", "", "Validate target api file")

	protoCmd.Flags().StringVar(&gogen.VarStringProto, "proto", "", "The proto path")
	protoCmd.Flags().StringVar(&gogen.VarStringOutput, "o", "", "The output path")
	protoCmd.Flags().StringVar(&gogen.VarStringServiceName, "serviceName", "", "The service name")
	protoCmd.Flags().StringVar(&gogen.VarStringStyle, "style", "go_zero", "The file name format style")
	protoCmd.Flags().StringVar(&gogen.VarStringModelName, "model", "", "The model name for generating e.g. user, "+
		"if it is empty, generate codes for all models in schema directory")
	protoCmd.Flags().IntVar(&gogen.VarIntSearchKeyNum, "searchKeyNum", 3, "The max number of search keys ")
	protoCmd.Flags().StringVar(&gogen.VarStringRpcName, "rpcName", "", "The rpc name in service context. e.g. CoreRpc")
	protoCmd.Flags().StringVar(&gogen.VarStringGrpcPbPackage, "grpcPackage", "", "The rpc name in service context. e.g. CoreRpc")

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
}
