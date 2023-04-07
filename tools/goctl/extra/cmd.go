package extra

import (
	"github.com/zeromicro/go-zero/tools/goctl/extra/i18n"
	"github.com/zeromicro/go-zero/tools/goctl/extra/initlogic"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

var (
	ExtraCmd = cobrax.NewCommand("extra")

	i18nCmd = cobrax.NewCommand("i18n", cobrax.WithRunE(i18n.Gen))

	initCmd = cobrax.NewCommand("init_code", cobrax.WithRunE(initlogic.Gen))
)

func init() {
	var (
		i18nCmdFlags = i18nCmd.Flags()
		initCmdFlags = initCmd.Flags()
	)

	i18nCmdFlags.StringVar(&i18n.VarStringTarget, "target")
	i18nCmdFlags.StringVar(&i18n.VarStringModelName, "model_name")
	i18nCmdFlags.StringVar(&i18n.VarStringModelNameZh, "model_name_zh")
	i18nCmdFlags.StringVar(&i18n.VarStringOutputDir, "output")

	initCmdFlags.StringVar(&initlogic.VarStringTarget, "target")
	initCmdFlags.StringVar(&initlogic.VarStringModelName, "model_name")
	initCmdFlags.StringVar(&initlogic.VarStringOutputPath, "output")

	ExtraCmd.AddCommand(i18nCmd)
	ExtraCmd.AddCommand(initCmd)
}
