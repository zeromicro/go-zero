package main

import (
	"fmt"
	"os"

	"zero/core/lang"
	"zero/tools/goctl/api/parser"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}

	p, err := parser.NewParser(os.Args[1])
	lang.Must(err)
	api, err := p.Parse()
	lang.Must(err)
	fmt.Println(api)
}
