package configgen

import (
	"fmt"
	"os"

	"zero/tools/modelctl/model"

	"github.com/urfave/cli"
)

var (
	configFileName = "config.json"
)

func ConfigCommand(_ *cli.Context) error {
	_, err := os.Stat(configFileName)
	if err == nil {
		return nil
	}
	file, err := os.OpenFile(configFileName, os.O_CREATE|os.O_WRONLY, model.ModeDirPerm)
	if err != nil {
		return err
	}
	file.WriteString(configTemplate)
	defer file.Close()
	fmt.Println("config json template generate done ... ")
	return nil
}
