package quickstart

import "github.com/spf13/cobra"

const (
	serviceTypeMono  = "mono"
	serviceTypeMicro = "micro"
)

var (
	varStringServiceType string

	// Cmd describes the command to run.
	Cmd = &cobra.Command{
		Use:   "quickstart",
		Short: "quickly start a project",
		RunE:  run,
	}
)

func init() {
	Cmd.Flags().StringVarP(&varStringServiceType,
		"service-type", "t", "mono",
		"specify the service type, supported values: [mono, micro]")
}
