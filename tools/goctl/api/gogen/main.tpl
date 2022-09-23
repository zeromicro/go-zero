{{.systemInfo}}
package main

import (
	"flag"
	"fmt"

	{{.importPackages}}
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")


func main() {
	flag.Parse()

	var consulConfig config.ConsulConfig
    conf.MustLoad(*configFile, &consulConfig)

    var c config.Config
    client, err := consulConfig.Consul.NewClient()
    logx.Must(err)
    consul.LoadYAMLConf(client, "{{.serviceName}}ApiConf", &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	err = consul.RegisterService(consulConfig.Consul)
    logx.Must(err)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
