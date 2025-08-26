package bug

import (
	"github.com/lerity-yao/go-zero/tools/cztctl/internal/cobrax"
	"github.com/spf13/cobra"
)

// Cmd describes a bug command.
var Cmd = cobrax.NewCommand("bug", cobrax.WithRunE(runE), cobrax.WithArgs(cobra.NoArgs))
