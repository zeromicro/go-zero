package generator

import (
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
