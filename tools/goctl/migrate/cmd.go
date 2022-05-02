package migrate

import "github.com/spf13/cobra"

var (
	boolVarVerbose   bool
	stringVarVersion string
	Cmd              = &cobra.Command{
		Use:   "migrate",
		Short: "migrate from tal-tech to zeromicro",
		Long: "migrate is a transition command to help users migrate their " +
			"projects from tal-tech to zeromicro version",
		RunE: migrate,
	}
)

func init() {
	Cmd.Flags().BoolVarP(&boolVarVerbose, "verbose", "v",
		false, "verbose enables extra logging")
	Cmd.Flags().StringVar(&stringVarVersion, "version", defaultMigrateVersion,
		"the target release version of github.com/zeromicro/go-zero to migrate")
}
