package bug

import "github.com/spf13/cobra"

// Cmd describes a bug command.
var Cmd = &cobra.Command{
	Use:   "bug",
	Short: "Report a bug",
	Args:  cobra.NoArgs,
	RunE:  runE,
}
