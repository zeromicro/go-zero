package migrate

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

var (
	boolVarVerbose   bool
	stringVarVersion string
	// Cmd describes a migrate command.
	Cmd = &cobra.Command{
		Use:   "migrate",
		Short: flags.Get("migrate.short"),
		Long:  flags.Get("migrate.long"),
		RunE:  migrate,
	}
)

func init() {
	migrateCmdFlags := Cmd.Flags()
	migrateCmdFlags.BoolVarP(&boolVarVerbose, "verbose", "v", false, flags.Get("migrate.verbose"))
	migrateCmdFlags.StringVar(&stringVarVersion, "version", defaultMigrateVersion, flags.Get("migrate.version"))
}
