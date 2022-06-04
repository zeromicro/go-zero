package parser

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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

func (s Services) validate(filename string, compatibleOpt ...bool) error {
	if len(s) == 0 {
		return errors.New("rpc service not found")
	}

	var compatible bool
	for _, c := range compatibleOpt {
		compatible = c
	}

	if compatible && len(s) > 1 {
		return errors.New("only one service expected")
	}

	name := filepath.Base(filename)
	for _, service := range s {
		for _, rpc := range service.RPC {
			if strings.Contains(rpc.RequestType, ".") {
				return fmt.Errorf("line %v:%v, request type must defined in %s",
					rpc.Position.Line,
					rpc.Position.Column, name)
			}
			if strings.Contains(rpc.ReturnsType, ".") {
				return fmt.Errorf("line %v:%v, returns type must defined in %s",
					rpc.Position.Line,
					rpc.Position.Column, name)
			}
		}
	}
	return nil
}
