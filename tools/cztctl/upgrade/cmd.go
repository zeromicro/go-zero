package upgrade

import "github.com/lerity-yao/go-zero/tools/cztctl/internal/cobrax"

// Cmd describes an upgrade command.
var Cmd = cobrax.NewCommand("upgrade", cobrax.WithRunE(upgrade))
