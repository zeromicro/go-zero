package parser

import "github.com/emicklei/proto"

func GetComment(comment *proto.Comment) string {
	if comment == nil {
		return ""
	}
	return "// " + comment.Message()
}
