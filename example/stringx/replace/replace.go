package main

import (
	"fmt"

	"github.com/3Rivers/go-zero/core/stringx"
)

func main() {
	replacer := stringx.NewReplacer(map[string]string{
		"日本":    "法国",
		"日本的首都": "东京",
		"东京":    "日本的首都",
	})
	fmt.Println(replacer.Replace("日本的首都是东京"))
}
