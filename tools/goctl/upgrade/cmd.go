package upgrade

import "github.com/spf13/cobra"

// Cmd describes the command to run.
var Cmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade goctl to latest version",
	RunE:  upgrade,
}
