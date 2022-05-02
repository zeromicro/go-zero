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
	Cmd = &cobra.Command{
		Use:   "api",
		Short: "generate api related files",
		RunE:  apigen.CreateApiTemplate,
	}

	dartCmd = &cobra.Command{
		Use:   "dart",
		Short: "generate dart files for provided api in api file",
		RunE:  dartgen.DartCommand,
	}

	docCmd = &cobra.Command{
		Use:   "doc",
		Short: "generate doc files",
		RunE:  docgen.DocCommand,
	}

	formatCmd = &cobra.Command{
		Use:   "format",
		Short: "format api files",
		RunE:  format.GoFormatApi,
	}

	goCmd = &cobra.Command{
		Use:   "go",
		Short: "generate go files for provided api in yaml file",
		RunE:  gogen.GoCommand,
	}

	newCmd = &cobra.Command{
		Use:     "new",
		Short:   "fast create api service",
		Example: "goctl api new [options] service-name",
		Args:    cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return new.CreateServiceCommand(args)
		},
	}

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "validate api file",
		RunE:  validate.GoValidateApi,
	}

	javaCmd = &cobra.Command{
		Use:   "java",
		Short: "generate java files for provided api in api file",
		RunE:  javagen.JavaCommand,
	}

	ktCmd = &cobra.Command{
		Use:   "kt",
		Short: "generate kotlin code for provided api file",
		RunE:  ktgen.KtCommand,
	}

	pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "custom file generator",
		RunE:  plugin.PluginCommand,
	}

	tsCmd = &cobra.Command{
		Use:   "ts",
		Short: "generate ts files for provided api in api file",
		RunE:  tsgen.TsCommand,
	}
)

func init() {
	Cmd.Flags().StringVar(&apigen.VarStringOutput, "o", "", "output a sample api file")
	Cmd.Flags().StringVar(&apigen.VarStringHome, "home", "", "the goctl home path of the"+
		" template, --home and --remote cannot be set at the same time, if they are, --remote has "+
		"higher priority")
	Cmd.Flags().StringVar(&apigen.VarStringRemote, "remote", "", "the remote git repo of the"+
		" template, --home and --remote cannot be set at the same time, if they are, --remote has higher"+
		" priority\n\tThe git repo directory must be consistent with the"+
		" https://github.com/zeromicro/go-zero-template directory structure")
	Cmd.Flags().StringVar(&apigen.VarStringBranch, "branch", "master", "the branch of the "+
		"remote repo, it does work with --remote")

	dartCmd.Flags().StringVar(&dartgen.VarStringDir, "dir", "", "the target dir")
	dartCmd.Flags().StringVar(&dartgen.VarStringAPI, "api", "", "the api file")
	dartCmd.Flags().BoolVar(&dartgen.VarStringLegacy, "legacy", false, "legacy generator for flutter v1")
	dartCmd.Flags().StringVar(&dartgen.VarStringHostname, "hostname", "", "hostname of the server")

	docCmd.Flags().StringVar(&docgen.VarStringDir, "dir", "", "the target dir")
	docCmd.Flags().StringVar(&docgen.VarStringOutput, "o", "", "the output markdown directory")

	formatCmd.Flags().StringVar(&format.VarStringDir, "dir", "", "the format target dir")
	formatCmd.Flags().BoolVar(&format.VarBoolIgnore, "iu", false, "ignore update")
	formatCmd.Flags().BoolVar(&format.VarBoolUseStdin, "stdin", false, "use stdin to input api"+
		" doc content, press \"ctrl + d\" to send EOF")
	formatCmd.Flags().BoolVar(&format.VarBoolSkipCheckDeclare, "declare", false, "use to skip check "+
		"api types already declare")

	goCmd.Flags().StringVar(&gogen.VarStringDir, "dir", "", "the target dir")
	goCmd.Flags().StringVar(&gogen.VarStringAPI, "api", "", "the api file")
	goCmd.Flags().StringVar(&gogen.VarStringHome, "home", "", "the goctl home path of "+
		"the template, --home and --remote cannot be set at the same time, if they are, --remote "+
		"has higher priority")
	goCmd.Flags().StringVar(&gogen.VarStringRemote, "remote", "", "the remote git repo "+
		"of the template, --home and --remote cannot be set at the same time, if they are, --remote"+
		" has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	goCmd.Flags().StringVar(&gogen.VarStringBranch, "branch", "master", "the branch of "+
		"the remote repo, it does work with --remote")
	goCmd.Flags().StringVar(&gogen.VarStringStyle, "style", "gozero", "the file naming format,"+
		" see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")

	javaCmd.Flags().StringVar(&javagen.VarStringDir, "dir", "", "the target dir")
	javaCmd.Flags().StringVar(&javagen.VarStringAPI, "api", "", "the api file")

	ktCmd.Flags().StringVar(&ktgen.VarStringDir, "dir", "", "the target dir")
	ktCmd.Flags().StringVar(&ktgen.VarStringAPI, "api", "", "the api file")
	ktCmd.Flags().StringVar(&ktgen.VarStringPKG, "pkg", "", "define package name for kotlin file")

	newCmd.Flags().StringVar(&new.VarStringHome, "home", "", "the goctl home path of "+
		"the template, --home and --remote cannot be set at the same time, if they are, --remote "+
		"has higher priority")
	newCmd.Flags().StringVar(&new.VarStringRemote, "remote", "", "the remote git repo "+
		"of the template, --home and --remote cannot be set at the same time, if they are, --remote"+
		" has higher priority\n\tThe git repo directory must be consistent with the "+
		"https://github.com/zeromicro/go-zero-template directory structure")
	newCmd.Flags().StringVar(&new.VarStringBranch, "branch", "master", "the branch of "+
		"the remote repo, it does work with --remote")
	newCmd.Flags().StringVar(&new.VarStringStyle, "style", "gozero", "the file naming format,"+
		" see [https://github.com/zeromicro/go-zero/blob/master/tools/goctl/config/readme.md]")

	pluginCmd.Flags().StringVarP(&plugin.VarStringPlugin, "plugin", "p", "", "the plugin file")
	pluginCmd.Flags().StringVar(&plugin.VarStringDir, "dir", "", "the target dir")
	pluginCmd.Flags().StringVar(&plugin.VarStringAPI, "api", "", "the api file")
	pluginCmd.Flags().StringVar(&plugin.VarStringStyle, "style", "",
		"the file naming format, see [https://github.com/zeromicro/go-zero/tree/master/tools/goctl/config/readme.md]")

	tsCmd.Flags().StringVar(&tsgen.VarStringDir, "dir", "", "the target dir")
	tsCmd.Flags().StringVar(&tsgen.VarStringAPI, "api", "", "the api file")
	tsCmd.Flags().StringVar(&tsgen.VarStringWebAPI, "webapi", "", "the web api file path")
	tsCmd.Flags().StringVar(&tsgen.VarStringCaller, "caller", "", "the web api caller")
	tsCmd.Flags().BoolVar(&tsgen.VarBoolUnWrap, "unwrap", false, "unwrap the webapi caller for import")

	validateCmd.Flags().StringVar(&validate.VarStringAPI, "api", "", "validate target api file")

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

}
