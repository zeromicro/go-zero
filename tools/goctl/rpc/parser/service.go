package parser

import "github.com/emicklei/proto"

// Service describes the rpc service, which is the relevant
// content after the translation of the proto file
type Service struct {
	*proto.Service
	RPC []*RPC
}
