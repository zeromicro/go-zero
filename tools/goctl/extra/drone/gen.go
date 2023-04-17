package drone

import "github.com/spf13/cobra"

var (
	// VarBoolDockerfile describes whether to generate dockerfile
	VarBoolDockerfile bool
)

func GenDrone(_ *cobra.Command, _ []string) error {
	return nil
}
