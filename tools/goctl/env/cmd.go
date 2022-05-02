package env

import "github.com/spf13/cobra"

var (
	sliceVarWriteValue []string
	boolVarForce       bool
	boolVarVerbose     bool
	boolVarInstall     bool

	// Cmd describes the command
	Cmd = &cobra.Command{
		Use:   "env",
		Short: "check or edit goctl environment",
		RunE:  write,
	}
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "goctl env installation",
		RunE:  install,
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "detect goctl env and dependency tools",
		RunE:  check,
	}
)

func init() {
	// The root command flags
	Cmd.Flags().StringSliceVarP(&sliceVarWriteValue,
		"write", "w", nil, "edit goctl environment")
	Cmd.PersistentFlags().BoolVarP(&boolVarForce,
		"force", "f", false,
		"silent installation of non-existent dependencies")
	Cmd.PersistentFlags().BoolVarP(&boolVarVerbose,
		"verbose", "v", false, "enable log output")

	// The sub-command flags
	checkCmd.Flags().BoolVarP(&boolVarInstall, "install", "i",
		false, "install dependencies if not found")

	// Add sub-command
	Cmd.AddCommand(installCmd)
	Cmd.AddCommand(checkCmd)
}
