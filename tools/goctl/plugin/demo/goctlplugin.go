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

	if plugin.Context != nil {
		fmt.Printf("context: %+v \n", plugin.Context)
	}

	fmt.Println("Enjoy anything you want.")
}
