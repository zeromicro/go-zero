package textgen

import (
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/textgen/i18n"
)

var (
	TextGenCmd = cobrax.NewCommand("textgen")

	i18nCmd = cobrax.NewCommand("i18n")

	initCmd = cobrax.NewCommand("init_text")
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

	initCmdFlags.StringVar(&i18n.VarStringTarget, "model_name")

	TextGenCmd.AddCommand(i18nCmd)
	TextGenCmd.AddCommand(initCmd)
}
