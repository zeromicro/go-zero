package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// Cmd describes a bug command.
	Cmd = cobrax.NewCommand("config")

	initCmd  = cobrax.NewCommand("init", cobrax.WithRunE(runConfigInit))
	cleanCmd = cobrax.NewCommand("clean", cobrax.WithRunE(runConfigClean))
)

func init() {
	Cmd.AddCommand(initCmd, cleanCmd)
}

func runConfigInit(*cobra.Command, []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFile, err := getConfigPath(wd)
	if err != nil {
		return err
	}
	if pathx.FileExists(cfgFile) {
		fmt.Printf("%s already exists, path: %s\n", configFile, cfgFile)
		return nil
	}

	err = os.WriteFile(cfgFile, defaultConfig, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("%s generated in %s\n", configFile, cfgFile)
	return nil
}

func runConfigClean(*cobra.Command, []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFile, err := getConfigPath(wd)
	if err != nil {
		return err
	}

	return pathx.RemoveIfExist(cfgFile)
}
