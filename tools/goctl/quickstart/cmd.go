package quickstart

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

const (
	serviceTypeMono  = "mono"
	serviceTypeMicro = "micro"
)

var (
	varStringServiceType string

	// Cmd describes the command to run.
	Cmd = &cobra.Command{
		Use:   "quickstart",
		Short: flags.Get("quickstart.short"),
		RunE:  run,
	}
)

func init() {
	Cmd.Flags().StringVarP(&varStringServiceType, "service-type", "t", "mono", flags.Get("quickstart.service-type"))
}
