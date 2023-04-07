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

	i18nCmdFlags.StringVarP(&i18n.VarStringTarget, "target", "t")
	i18nCmdFlags.StringVarP(&i18n.VarStringModelName, "model_name", "m")
	i18nCmdFlags.StringVarP(&i18n.VarStringModelNameZh, "model_name_zh", "z")
	i18nCmdFlags.StringVarP(&i18n.VarStringOutputDir, "output", "o")

	initCmdFlags.StringVarP(&initlogic.VarStringTarget, "target", "t")
	initCmdFlags.StringVarP(&initlogic.VarStringModelName, "model_name", "m")
	initCmdFlags.StringVarP(&initlogic.VarStringOutputPath, "output", "o")

	ExtraCmd.AddCommand(i18nCmd)
	ExtraCmd.AddCommand(initCmd)
}
