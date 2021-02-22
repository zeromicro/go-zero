package parser

import "github.com/emicklei/proto"

// Message embeds proto.Message
type Message struct {
	*proto.Message
}
