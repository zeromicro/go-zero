package main

import (
	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/cmd"
)

const codeFailure = 1

var commands = []cli.Command{}

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
	// cli.BashCompletionFlag = cli.BoolFlag{
	// 	Name:   completion.BashCompletionFlag,
	// 	Hidden: true,
	// }
	// app := cli.NewApp()
	// app.EnableBashCompletion = true
	// app.Usage = "a cli tool to generate code"
	// app.Version = fmt.Sprintf("%s %s/%s", version.BuildVersion, runtime.GOOS, runtime.GOARCH)
	// app.Commands = commands
	//
	// // cli already print error messages.
	// if err := app.Run(os.Args); err != nil {
	// 	fmt.Println(aurora.Red(err.Error()))
	// 	os.Exit(codeFailure)
	// }
}
