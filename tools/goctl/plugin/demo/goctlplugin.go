package main

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/plugin"
)

func main() {
	plugin, err := plugin.NewPlugin()
	if err != nil {
		panic(err)
	}

	if plugin.Api != nil {
		fmt.Printf("api: %+v \n", plugin.Api)
	}
	fmt.Printf("dir: %s \n", plugin.Dir)
	fmt.Println("Enjoy anything you want.")
}
