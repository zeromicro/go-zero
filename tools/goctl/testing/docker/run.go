package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

const tag = "goctl:latest"

func Run(filePath string) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println("docker must be started")
		os.Exit(-1)
	}
	// first build image from Dockerfile.
	// it depend on goctl compiled file.
	image := newImage(cli)
	if err := image.build("./", []string{tag}); err != nil {
		panic(err)
	}
	// parse docker compose yaml file.
	// start the container and execute goctl command.
	compose := newCompose(cli)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	cfg, err := decodeConfig(data)
	if err != nil {
		panic(err)
	}
	dir, _ := filepath.Abs(filepath.Dir(filePath))
	for name, v := range cfg.Services {
		if ct, err := newContainer(cli, name, v, dir); err != nil {
			fmt.Printf("new container: %v error: %+v\n", name, err)
			os.Exit(-1)
		} else {
			compose.add(ct)
		}
	}
	if err := compose.Run(); err != nil {
		fmt.Printf("compose run failed err: %v\n", err)
		os.Exit(-1)
	}
}
