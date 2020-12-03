package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

var (
	specStr    = flag.String("spec", "", "the spec file")
	contextStr = flag.String("context", "", "the context")
)

func main() {
	flag.Parse()

	var api spec.ApiSpec
	err := json.Unmarshal([]byte(*specStr), &api)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", api)

	var context ctx.ProjectContext
	err = json.Unmarshal([]byte(*contextStr), &context)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", context)
	fmt.Print("do any thing you can.")
}
