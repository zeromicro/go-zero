package completion

import "github.com/spf13/cobra"

var (
	varStringName string
	Cmd           = &cobra.Command{
		Use:   "completion",
		Short: "generation completion script, it only works for unix-like OS",
		RunE:  completion,
	}
)

func init() {
	Cmd.Flags().StringVarP(&varStringName, "name", "n", "goctl_autocomplete", "the filename of auto complete script, default is [goctl_autocomplete]")
}
