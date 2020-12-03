package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

var (
	specFile    = flag.String("spec", "", "the spec file")
	contextFile = flag.String("context", "", "the context")
)

func main() {
	flag.Parse()

	var api spec.ApiSpec
	content, err := ioutil.ReadFile(*specFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &api)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", api)

	var context ctx.ProjectContext
	content, err = ioutil.ReadFile(*contextFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &context)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", context)
	fmt.Println("Enjoy anything you can.")
}
