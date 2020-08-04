package main

import (
	"zero/core/discov"
	"zero/core/lang"
	"zero/rpcx/internal"
)

func main() {
	cli, err := internal.NewDiscovClient(discov.EtcdConf{
		Hosts: []string{"localhost:2379"},
		Key:   "rpcx",
	})
	lang.Must(err)

	cli.Next()
}
