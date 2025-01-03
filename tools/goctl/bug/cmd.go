package bug

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

// Cmd describes a bug command.
var Cmd = cobrax.NewCommand("bug", cobrax.WithRunE(runE), cobrax.WithArgs(cobra.NoArgs))
