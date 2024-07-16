package web

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server"
)

var (
	portFlag int

	// Cmd describes a rpc command.
	Cmd = cobrax.NewCommand("web", cobrax.WithRunE(
		func(command *cobra.Command, strings []string) error {
			return server.Run(portFlag)
		}),
	)
)

func init() {
	Cmd.Flags().IntVarP(&portFlag, "port", "p", server.DefaultPort, "server port")
}
