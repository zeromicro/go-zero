package generator

import (
	"os/exec"

	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

type defaultGenerator struct {
	log console.Console
}

func NewDefaultGenerator() *defaultGenerator {
	log := console.NewColorConsole()
	return &defaultGenerator{
		log: log,
	}
}

func (g *defaultGenerator) Prepare() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return err
	}

	_, err = exec.LookPath("protoc")
	if err != nil {
		return err
	}

	_, err = exec.LookPath("protoc-gen-go")
	return err
}
