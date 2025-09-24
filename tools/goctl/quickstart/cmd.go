package quickstart

import "github.com/dropship-dev/go-zero/tools/goctl/internal/cobrax"

const (
	serviceTypeMono  = "mono"
	serviceTypeMicro = "micro"
)

var (
	varStringServiceType string

	// Cmd describes the command to run.
	Cmd = cobrax.NewCommand("quickstart", cobrax.WithRunE(run))
)

func init() {
	Cmd.Flags().StringVarPWithDefaultValue(&varStringServiceType, "service-type", "t", "mono")
}
