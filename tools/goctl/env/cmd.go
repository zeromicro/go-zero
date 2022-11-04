package env

import "github.com/spf13/cobra"

var (
	sliceVarWriteValue []string
	boolVarForce       bool
	boolVarVerbose     bool
	boolVarInstall     bool

	// Cmd describes an env command.
	Cmd = &cobra.Command{
		Use:   "env",
		Short: "Check or edit goctl environment",
		RunE:  write,
	}
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Goctl env installation",
		RunE:  install,
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Detect goctl env and dependency tools",
		RunE:  check,
	}
)

func init() {
	// The root command flags
	Cmd.Flags().StringSliceVarP(&sliceVarWriteValue,
		"write", "w", nil, "Edit goctl environment")
	Cmd.PersistentFlags().BoolVarP(&boolVarForce,
		"force", "f", false,
		"Silent installation of non-existent dependencies")
	Cmd.PersistentFlags().BoolVarP(&boolVarVerbose,
		"verbose", "v", false, "Enable log output")

	// The sub-command flags
	checkCmd.Flags().BoolVarP(&boolVarInstall, "install", "i",
		false, "Install dependencies if not found")

	// Add sub-command
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(checkCmd)
}
