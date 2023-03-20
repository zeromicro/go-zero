package env

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

var (
	sliceVarWriteValue []string
	boolVarForce       bool
	boolVarVerbose     bool
	boolVarInstall     bool

	// Cmd describes an env command.
	Cmd = &cobra.Command{
		Use:   "env",
		Short: flags.Get("env.short"),
		RunE:  write,
	}
	installCmd = &cobra.Command{
		Use:   "install",
		Short: flags.Get("env.install.short"),
		RunE:  install,
	}
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: flags.Get("env.check.short"),
		RunE:  check,
	}
)

func init() {
	// The root command flags
	Cmd.Flags().StringSliceVarP(&sliceVarWriteValue, "write", "w", nil, flags.Get("env.write"))
	Cmd.PersistentFlags().BoolVarP(&boolVarForce, "force", "f", false, flags.Get("env.force"))
	Cmd.PersistentFlags().BoolVarP(&boolVarVerbose, "verbose", "v", false, flags.Get("env.verbose"))

	// The sub-command flags
	checkCmd.Flags().BoolVarP(&boolVarInstall, "install", "i", false, flags.Get("env.check.install"))

	// Add sub-command
	Cmd.AddCommand(checkCmd, installCmd)
}
