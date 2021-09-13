package generator

import (
	"os/exec"

	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

// DefaultGenerator defines the environment needs of rpc service generation
type DefaultGenerator struct {
	log console.Console
}

// just test interface implement
var _ Generator = (*DefaultGenerator)(nil)

// NewDefaultGenerator returns an instance of DefaultGenerator
func NewDefaultGenerator() Generator {
	log := console.NewColorConsole()
	return &DefaultGenerator{
		log: log,
	}
}

// Prepare provides environment detection generated by rpc service,
// including go environment, protoc, whether protoc-gen-go is installed or not
func (g *DefaultGenerator) Prepare() error {
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
