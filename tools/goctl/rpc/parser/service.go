package parser

import "github.com/emicklei/proto"

type Service struct {
	*proto.Service
	RPC []*RPC
}
