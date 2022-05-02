package bug

import "github.com/spf13/cobra"

// Cmd describes the command to run.
var Cmd = &cobra.Command{
	Use:   "bug",
	Short: "report a bug",
	Args:  cobra.NoArgs,
	RunE:  runE,
}
