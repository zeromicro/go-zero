package migrate

import "github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"

var (
	boolVarVerbose bool
	zeroVersion    string
	toolVersion    string
	// Cmd describes a migrate command.
	Cmd = cobrax.NewCommand("migrate", cobrax.WithRunE(migrate))
)

func init() {
	migrateCmdFlags := Cmd.Flags()
	migrateCmdFlags.BoolVarP(&boolVarVerbose, "verbose", "v")
	migrateCmdFlags.StringVar(&zeroVersion, "zero_version")
	migrateCmdFlags.StringVar(&toolVersion, "tool_version")
}
