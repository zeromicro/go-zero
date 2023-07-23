package env

import "github.com/spf13/cobra"

func install(_ *cobra.Command, _ []string) error {
	return Prepare(true, boolVarForce, boolVarVerbose)
}
