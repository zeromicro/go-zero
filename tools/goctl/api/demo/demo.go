package main

import (
	"flag"

	"zero/core/conf"
	"zero/ngin"
	"zero/tools/goctl/api/demo/config"
	"zero/tools/goctl/api/demo/handler"
)

var configFile = flag.String("f", "etc/user.json", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	engine := ngin.MustNewEngine(c.NgConf)
	defer engine.Stop()

	handler.RegisterHandlers(engine)
	engine.Start()
}
