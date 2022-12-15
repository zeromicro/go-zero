package parser

import (
	"strings"

	"github.com/emicklei/proto"
)

// GetComment returns content with prefix //
func GetComment(comment *proto.Comment) string {
	if comment == nil || strings.Contains(comment.Message(), "group") {
		return ""
	}
	return "// " + strings.TrimSpace(comment.Message())
}
