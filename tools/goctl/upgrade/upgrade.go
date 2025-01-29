package upgrade

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

// upgrade gets the latest goctl by
// go install github.com/zeromicro/go-zero/tools/goctl@latest
func upgrade(_ *cobra.Command, _ []string) error {
	cmd := `go install github.com/zeromicro/go-zero/tools/goctl@latest`
	if runtime.GOOS == "windows" {
		cmd = `go install github.com/zeromicro/go-zero/tools/goctl@latest`
	}
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
