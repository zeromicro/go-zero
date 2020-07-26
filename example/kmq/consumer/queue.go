package main

import (
	"fmt"

	"zero/core/conf"
	"zero/kq"
)

func main() {
	var c kq.KqConf
	conf.MustLoad("config.json", &c)

	q := kq.MustNewQueue(c, kq.WithHandle(func(k, v string) error {
		fmt.Printf("=> %s\n", v)
		return nil
	}))
	defer q.Stop()
	q.Start()
}
