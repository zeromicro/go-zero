package migrate

import (
	"github.com/spf13/cobra"
)

const defaultMigrateVersion = "v1.5.1"

func migrate(_ *cobra.Command, _ []string) error {
	if len(zeroVersion) == 0 {
		zeroVersion = defaultMigrateVersion
	}
	err := editMod(zeroVersion, toolVersion, boolVarVerbose)
	if err != nil {
		return err
	}

	err = tidy(boolVarVerbose)
	if err != nil {
		return err
	}

	return nil
}
