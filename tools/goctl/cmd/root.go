package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/api"
	"github.com/zeromicro/go-zero/tools/goctl/bug"
	"github.com/zeromicro/go-zero/tools/goctl/docker"
	"github.com/zeromicro/go-zero/tools/goctl/env"
	"github.com/zeromicro/go-zero/tools/goctl/kube"
	"github.com/zeromicro/go-zero/tools/goctl/migrate"
	"github.com/zeromicro/go-zero/tools/goctl/model"
	"github.com/zeromicro/go-zero/tools/goctl/rpc"
	"github.com/zeromicro/go-zero/tools/goctl/tpl"
	"github.com/zeromicro/go-zero/tools/goctl/upgrade"
)

const codeFailure = 1

var rootCmd = &cobra.Command{
	Use:   "goctl",
	Short: "a cli tool to generate go-zero code",
	Long:  "a cli tool to generate api, zrpc, model code",
}

// Execute executes the given command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(aurora.Red(err.Error()))
		os.Exit(codeFailure)
	}
}

func init() {
	rootCmd.AddCommand(api.Cmd)
	rootCmd.AddCommand(bug.Cmd)
	// rootCmd.AddCommand(completion.Cmd)
	rootCmd.AddCommand(docker.Cmd)
	rootCmd.AddCommand(kube.Cmd)
	rootCmd.AddCommand(env.Cmd)
	rootCmd.AddCommand(model.Cmd)
	rootCmd.AddCommand(migrate.Cmd)
	rootCmd.AddCommand(rpc.Cmd)
	rootCmd.AddCommand(tpl.Cmd)
	rootCmd.AddCommand(upgrade.Cmd)
}
