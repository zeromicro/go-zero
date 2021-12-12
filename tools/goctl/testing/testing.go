package main

import (
	"os/exec"

	"github.com/tal-tech/go-zero/tools/goctl/testing/docker"
)

func main() {
	if err := compile(); err != nil {
		panic(err)
	}
	defer del()
	docker.Run("./docker-compose.yaml")
}

func compile() error {
	if err := exec.Command("make", "-C", "..", "linux").Run(); err != nil {
		return err
	}
	return exec.Command("mv", "../goctl-linux", "./goctl").Run()
}

func del() error {
	return exec.Command("rm", "-f", "./goctl").Run()
}
