package parser

import (
	"errors"

	"github.com/emicklei/proto"
)

type (
	// Services is a slice of Service.
	Services []Service

	// Service describes the rpc service, which is the relevant
	// content after the translation of the proto file
	Service struct {
		*proto.Service
		RPC []*RPC
	}
)

func (s Services) validate(filename string, multipleOpt ...bool) error {
	if len(s) == 0 {
		// return errors.New("rpc service not found")

		// message only proto file, no service defined
		return nil
	}

	var multiple bool
	for _, c := range multipleOpt {
		multiple = c
	}

	if !multiple && len(s) > 1 {
		return errors.New("only one service expected")
	}
	return nil
}
