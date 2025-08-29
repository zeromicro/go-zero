package main

import (
	"flag"
	"fmt"
	"github.com/lerity-yao/go-zero/tools/cztctl/bbbbf/internal/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/bbbbf/internal/handler"
	"github.com/lerity-yao/go-zero/tools/cztctl/bbbbf/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/cztAdminPc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	handler.RegisterHandlers(serviceGroup, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	serviceGroup.Start()
}
