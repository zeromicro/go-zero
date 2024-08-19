package quickstart

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const baseDir = "greet"

var (
	log        = console.NewColorConsole(true)
	projectDir string
)

func cleanWorkSpace(projectDir string) {
	var command string
	var breakeState bool
	fmt.Printf("Detected that the %q already exists, do you clean up?"+
		" [y: YES, n: NO]: ", projectDir)

	for {
		fmt.Scanln(&command)
		switch command {
		case "y":
			log.Debug("Clean workspace...")
			_ = os.RemoveAll(projectDir)
			breakeState = true
			break
		case "n":
			log.Error("User canceled")
			os.Exit(1)
		default:
			fmt.Println("Invalid command, try again...")
		}

		if breakeState {
			break
		}
	}
}

func initProject() {
	wd, err := os.Getwd()
	logx.Must(err)

	projectDir = filepath.Join(wd, baseDir)
	if exists := pathx.FileExists(projectDir); exists {
		cleanWorkSpace(projectDir)
	}

	log.Must(pathx.MkdirIfNotExist(projectDir))
	_, err = ctx.Prepare(projectDir)
	logx.Must(err)
}

func run(_ *cobra.Command, _ []string) error {
	initProject()
	switch varStringServiceType {
	case serviceTypeMono:
		newMonoService(false).start()
	case serviceTypeMicro:
		newMicroService().start()
	default:
		return fmt.Errorf("invalid service type, expected %s | %s",
			serviceTypeMono, serviceTypeMicro)
	}
	return nil
}
