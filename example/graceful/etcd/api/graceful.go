package main

import (
	"flag"

	"zero/core/conf"
	"zero/example/graceful/etcd/api/config"
	"zero/example/graceful/etcd/api/handler"
	"zero/example/graceful/etcd/api/svc"
	"zero/ngin"
	"zero/rpcx"
)

var configFile = flag.String("f", "etc/graceful-api.json", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	client := rpcx.MustNewClient(c.Rpc)
	ctx := &svc.ServiceContext{
		Client: client,
	}

	engine := ngin.MustNewEngine(c.NgConf)
	defer engine.Stop()

	handler.RegisterHandlers(engine, ctx)
	engine.Start()
}
