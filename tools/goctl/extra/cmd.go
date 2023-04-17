package extra

import (
	"github.com/zeromicro/go-zero/tools/goctl/extra/drone"
	"github.com/zeromicro/go-zero/tools/goctl/extra/ent/template"
	"github.com/zeromicro/go-zero/tools/goctl/extra/i18n"
	"github.com/zeromicro/go-zero/tools/goctl/extra/initlogic"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

var (
	ExtraCmd = cobrax.NewCommand("extra")

	i18nCmd = cobrax.NewCommand("i18n", cobrax.WithRunE(i18n.Gen))

	initCmd = cobrax.NewCommand("init_code", cobrax.WithRunE(initlogic.Gen))

	entCmd = cobrax.NewCommand("ent")

	templateCmd = cobrax.NewCommand("template", cobrax.WithRunE(template.GenTemplate))

	droneCmd = cobrax.NewCommand("drone", cobrax.WithRunE(drone.GenDrone))
)

func init() {
	var (
		i18nCmdFlags     = i18nCmd.Flags()
		initCmdFlags     = initCmd.Flags()
		templateCmdFlags = templateCmd.Flags()
		droneCmdFlags    = droneCmd.Flags()
	)

	i18nCmdFlags.StringVarP(&i18n.VarStringTarget, "target", "t")
	i18nCmdFlags.StringVarP(&i18n.VarStringModelName, "model_name", "m")
	i18nCmdFlags.StringVarP(&i18n.VarStringModelNameZh, "model_name_zh", "z")
	i18nCmdFlags.StringVarP(&i18n.VarStringOutputDir, "output", "o")

	initCmdFlags.StringVarP(&initlogic.VarStringTarget, "target", "t")
	initCmdFlags.StringVarP(&initlogic.VarStringModelName, "model_name", "m")
	initCmdFlags.StringVarP(&initlogic.VarStringOutputPath, "output", "o")

	templateCmdFlags.StringVarP(&template.VarStringDir, "dir", "d")
	templateCmdFlags.StringVarP(&template.VarStringAdd, "add", "a")
	templateCmdFlags.BoolVarP(&template.VarBoolList, "list", "l")
	templateCmdFlags.BoolVarP(&template.VarBoolUpdate, "update", "u")

	droneCmdFlags.BoolVarP(&drone.VarBoolDockerfile, "dockerfile", "d")

	ExtraCmd.AddCommand(i18nCmd)
	ExtraCmd.AddCommand(initCmd)
	entCmd.AddCommand(templateCmd)
	ExtraCmd.AddCommand(entCmd)
	ExtraCmd.AddCommand(droneCmd)
}
