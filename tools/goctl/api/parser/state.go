package parser

import "github.com/tal-tech/go-zero/tools/goctl/api/spec"

type state interface {
	process(api *spec.ApiSpec) (state, error)
}
