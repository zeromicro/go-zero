package main

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/stringx"
)

func main() {
	filter := stringx.NewTrie([]string{
		"AV演员",
		"苍井空",
		"AV",
		"日本AV女优",
		"AV演员色情",
	}, stringx.WithMask('?'))
	safe, keywords, found := filter.Filter("日本AV演员兼电视、电影演员。苍井空AV女优是xx出道, 日本AV女优们最精彩的表演是AV演员色情表演")
	fmt.Println(safe)
	fmt.Println(keywords)
	fmt.Println(found)
}
