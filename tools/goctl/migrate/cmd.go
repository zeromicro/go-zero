package migrate

import "github.com/spf13/cobra"

var (
	boolVarVerbose   bool
	stringVarVersion string
	// Cmd describes a migrate command.
	Cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate from tal-tech to zeromicro",
		Long: "Migrate is a transition command to help users migrate their " +
			"projects from tal-tech to zeromicro version",
		RunE: migrate,
	}
)

func init() {
	Cmd.Flags().BoolVarP(&boolVarVerbose, "verbose", "v",
		false, "Verbose enables extra logging")
	Cmd.Flags().StringVar(&stringVarVersion, "version", defaultMigrateVersion,
		"The target release version of github.com/zeromicro/go-zero to migrate")
}
