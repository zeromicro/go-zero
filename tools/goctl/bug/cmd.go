package bug

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

// Cmd describes a bug command.
var Cmd = &cobra.Command{
	Use:   "bug",
	Short: flags.Get("bug.short"),
	Args:  cobra.NoArgs,
	RunE:  runE,
}
