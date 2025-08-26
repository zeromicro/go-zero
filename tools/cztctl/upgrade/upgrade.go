package upgrade

import (
	"fmt"
	"runtime"

	"github.com/lerity-yao/go-zero/tools/cztctl/rpc/execx"
	"github.com/spf13/cobra"
)

// upgrade gets the latest cztctl by
// go install github.com/lerity-yao/go-zero/tools/cztctl@latest
func upgrade(_ *cobra.Command, _ []string) error {
	cmd := `go install github.com/lerity-yao/go-zero/tools/cztctl@latest`
	if runtime.GOOS == "windows" {
		cmd = `go install github.com/lerity-yao/go-zero/tools/cztctl@latest`
	}
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
