package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
	"github.com/zeromicro/go-zero/tools/goctl/test"
)

func Test_executeTest(t *testing.T) {
	executor := test.NewExecutor[string, string]()
	commands := append([]*cobra.Command{rootCmd}, getCommandsRecursively(rootCmd)...)
	for _, command := range commands {
		commandKey := getCommandName(command)
		fmt.Println(">>>>>: " + commandKey)
		command.Flags().VisitAll(func(flag *pflag.Flag) {
			name := flag.Name
			fmt.Println(commandKey + "." + name)
		})
		fmt.Println("<<<<<")
	}
	executor.Run(t, func(s string) string {
		return flags.Get(s)
	})
}

func getCommandName(cmd *cobra.Command) string {
	if cmd.HasParent() {
		return getCommandName(cmd.Parent()) + "." + cmd.Name()
	}
	return cmd.Name()
}
