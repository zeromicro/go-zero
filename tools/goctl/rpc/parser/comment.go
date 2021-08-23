package parser

import "github.com/emicklei/proto"

// GetComment returns content with prefix //
func GetComment(comment *proto.Comment) string {
	if comment == nil {
		return ""
	}
	return "// " + comment.Message()
}
