package migrate

import "github.com/lerity-yao/go-zero/tools/cztctl/internal/cobrax"

var (
	boolVarVerbose   bool
	stringVarVersion string
	// Cmd describes a migrate command.
	Cmd = cobrax.NewCommand("migrate", cobrax.WithRunE(migrate))
)

func init() {
	migrateCmdFlags := Cmd.Flags()
	migrateCmdFlags.BoolVarP(&boolVarVerbose, "verbose", "v")
	migrateCmdFlags.StringVarWithDefaultValue(&stringVarVersion, "version", defaultMigrateVersion)
}
