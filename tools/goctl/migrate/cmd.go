package migrate

import "github.com/spf13/cobra"

var (
	boolVarVerbose bool
	zeroVersion    string
	toolVersion    string
	// Cmd describes a migrate command.
	Cmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate from zeromicro to simple admin tools",
		Long: "Migrate is a transition command to help users migrate their " +
			"projects from tal-tech to zeromicro version",
		RunE: migrate,
	}
)

func init() {
	Cmd.Flags().BoolVarP(&boolVarVerbose, "verbose", "v",
		false, "Verbose enables extra logging")
	Cmd.Flags().StringVar(&zeroVersion, "zero-version", defaultMigrateVersion,
		"The target release version of github.com/zeromicro/go-zero to migrate")
	Cmd.Flags().StringVar(&toolVersion, "tool-version", "v0.0.6",
		"The target release version of github.com/suyuan32/simple-admin-tools to migrate")
}
