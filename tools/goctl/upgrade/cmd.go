package upgrade

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

// Cmd describes an upgrade command.
var Cmd = &cobra.Command{
	Use:   "upgrade",
	Short: flags.Get("upgrade.short"),
	RunE:  upgrade,
}
