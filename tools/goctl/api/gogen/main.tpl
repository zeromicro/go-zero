package main

import (
	"flag"
	"fmt"
	
	"github.com/num30/config"
	{{.importPackages}}
)

func main() {
	var c cfg.Config

	// search for config file in the etc folder without caring for the file extension
	// Prority: flags > env > config file
	err := config.NewConfReader("config").WithSearchDirs("etc", ".").Read(&c)
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}


	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
