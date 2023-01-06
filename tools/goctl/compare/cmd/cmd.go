package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/compare/testdata"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

var rootCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare the goctl commands generated results between urfave and cobra",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		testdata.MustRun(dir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		console.Error("%+v", err)
	}
}
