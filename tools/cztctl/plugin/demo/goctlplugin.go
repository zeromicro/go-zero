package main

import (
	"fmt"

	"github.com/lerity-yao/go-zero/tools/cztctl/plugin"
)

func main() {
	plugin, err := plugin.NewPlugin()
	if err != nil {
		panic(err)
	}

	if plugin.Api != nil {
		fmt.Printf("api: %+v \n", plugin.Api)
	}
	fmt.Println("Enjoy anything you want.")
}
