package main

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/stringx"
)

func main() {
	replacer := stringx.NewReplacer(map[string]string{
		"PHP": "PPT",
		"世界上": "吹牛",
	})
	fmt.Println(replacer.Replace("PHP是世界上最好的语言！"))
}
