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
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

var (
	// Cmd describes an api command.
	Cmd = &cobra.Command{
		Use:   "api",
		Short: flags.Get("api.short"),
		RunE:  apigen.CreateApiTemplate,
	}

	dartCmd = &cobra.Command{
		Use:   "dart",
		Short: flags.Get("api.dart.short"),
		RunE:  dartgen.DartCommand,
	}

	docCmd = &cobra.Command{
		Use:   "doc",
		Short: flags.Get("api.doc.short"),
		RunE:  docgen.DocCommand,
	}

	formatCmd = &cobra.Command{
		Use:   "format",
		Short: flags.Get("api.format.short"),
		RunE:  format.GoFormatApi,
	}

	goCmd = &cobra.Command{
		Use:   "go",
		Short: flags.Get("api.go.short"),
		RunE:  gogen.GoCommand,
	}

	newCmd = &cobra.Command{
		Use:     "new",
		Short:   flags.Get("api.new.short"),
		Example: flags.Get("api.new.example"),
		Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			return new.CreateServiceCommand(args)
		},
	}

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: flags.Get("api.validate.short"),
		RunE:  validate.GoValidateApi,
	}

	javaCmd = &cobra.Command{
		Use:    "java",
		Short:  flags.Get("api.java.short"),
		Hidden: true,
		RunE:   javagen.JavaCommand,
	}

	ktCmd = &cobra.Command{
		Use:   "kt",
		Short: flags.Get("api.kt.short"),
		RunE:  ktgen.KtCommand,
	}

	pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: flags.Get("api.plugin.short"),
		RunE:  plugin.PluginCommand,
	}

	tsCmd = &cobra.Command{
		Use:   "ts",
		Short: flags.Get("api.ts.short"),
		RunE:  tsgen.TsCommand,
	}
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
	)

	apiCmdFlags.StringVar(&apigen.VarStringOutput, "o", "", flags.Get("api.o"))
	apiCmdFlags.StringVar(&apigen.VarStringHome, "home", "", flags.Get("api.home"))
	apiCmdFlags.StringVar(&apigen.VarStringRemote, "remote", "", flags.Get("api.remote"))
	apiCmdFlags.StringVar(&apigen.VarStringBranch, "branch", "", flags.Get("api.branch"))

	dartCmdFlags.StringVar(&dartgen.VarStringDir, "dir", "", flags.Get("api.dart.dir"))
	dartCmdFlags.StringVar(&dartgen.VarStringAPI, "api", "", flags.Get("api.dart.api"))
	dartCmdFlags.BoolVar(&dartgen.VarStringLegacy, "legacy", false, flags.Get("api.dart.legacy"))
	dartCmdFlags.StringVar(&dartgen.VarStringHostname, "hostname", "", flags.Get("api.dart.hostname"))
	dartCmdFlags.StringVar(&dartgen.VarStringScheme, "scheme", "", flags.Get("api.dart.scheme"))

	docCmdFlags.StringVar(&docgen.VarStringDir, "dir", "", flags.Get("api.doc.dir"))
	docCmdFlags.StringVar(&docgen.VarStringOutput, "o", "", flags.Get("api.doc.o"))

	formatCmdFlags.StringVar(&format.VarStringDir, "dir", "", flags.Get("api.format.dir"))
	formatCmdFlags.BoolVar(&format.VarBoolIgnore, "iu", false, flags.Get("api.format.iu"))
	formatCmdFlags.BoolVar(&format.VarBoolUseStdin, "stdin", false, flags.Get("api.format.stdin"))
	formatCmdFlags.BoolVar(&format.VarBoolSkipCheckDeclare, "declare", false, flags.Get("api.format.declare"))

	goCmdFlags.StringVar(&gogen.VarStringDir, "dir", "", flags.Get("api.go.dir"))
	goCmdFlags.StringVar(&gogen.VarStringAPI, "api", "", flags.Get("api.go.api"))
	goCmdFlags.StringVar(&gogen.VarStringHome, "home", "", flags.Get("api.go.home"))
	goCmdFlags.StringVar(&gogen.VarStringRemote, "remote", "", flags.Get("api.go.remote"))
	goCmdFlags.StringVar(&gogen.VarStringBranch, "branch", "", flags.Get("api.go.branch"))
	goCmdFlags.StringVar(&gogen.VarStringStyle, "style", config.DefaultFormat, flags.Get("api.go.style"))

	javaCmdFlags.StringVar(&javagen.VarStringDir, "dir", "", flags.Get("api.java.dir"))
	javaCmdFlags.StringVar(&javagen.VarStringAPI, "api", "", flags.Get("api.java.api"))

	ktCmdFlags.StringVar(&ktgen.VarStringDir, "dir", "", flags.Get("api.kt.dir"))
	ktCmdFlags.StringVar(&ktgen.VarStringAPI, "api", "", flags.Get("api.kt.api"))
	ktCmdFlags.StringVar(&ktgen.VarStringPKG, "pkg", "", flags.Get("api.kt.pkg"))

	newCmdFlags.StringVar(&new.VarStringHome, "home", "", flags.Get("api.new.home"))
	newCmdFlags.StringVar(&new.VarStringRemote, "remote", "", flags.Get("api.new.remote"))
	newCmdFlags.StringVar(&new.VarStringBranch, "branch", "", flags.Get("api.new.branch"))
	newCmdFlags.StringVar(&new.VarStringStyle, "style", config.DefaultFormat, flags.Get("api.new.style"))

	pluginCmdFlags.StringVarP(&plugin.VarStringPlugin, "plugin", "p", "", flags.Get("api.plugin.plugin"))
	pluginCmdFlags.StringVar(&plugin.VarStringDir, "dir", "", flags.Get("api.plugin.dir"))
	pluginCmdFlags.StringVar(&plugin.VarStringAPI, "api", "", flags.Get("api.plugin.api"))
	pluginCmdFlags.StringVar(&plugin.VarStringStyle, "style", "", flags.Get("api.plugin.style"))

	tsCmdFlags.StringVar(&tsgen.VarStringDir, "dir", "", flags.Get("api.ts.dir"))
	tsCmdFlags.StringVar(&tsgen.VarStringAPI, "api", "", flags.Get("api.ts.api"))
	tsCmdFlags.StringVar(&tsgen.VarStringCaller, "caller", "", flags.Get("api.ts.caller"))
	tsCmdFlags.BoolVar(&tsgen.VarBoolUnWrap, "unwrap", false, flags.Get("api.ts.unwrap"))

	validateCmdFlags.StringVar(&validate.VarStringAPI, "api", "", flags.Get("api.validate.api"))

	// Add sub-commands
	Cmd.AddCommand(dartCmd, docCmd, formatCmd, goCmd, javaCmd, ktCmd, newCmd, pluginCmd, tsCmd, validateCmd)
}
