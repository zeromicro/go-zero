package protox

import (
	"fmt"
	"strings"

	"github.com/emicklei/proto"
)

func FindBeginEndOfService(service, serviceName string) (begin, mid, end int) {
	beginIndex := strings.Index(service, serviceName)
	begin, end = -1, -1
	if beginIndex > 0 {
		for i := beginIndex; i < len(service); i++ {
			if service[i] == '}' {
				end = i
				break
			} else if service[i] == '{' {
				mid = i
			}
		}

		for i := beginIndex; i >= 0; i-- {
			if service[i] == 's' {
				begin = i
				break
			}
		}
	}
	return begin, mid, end
}

var ProtoField *ProtoFieldData

type ProtoFieldData struct {
	Name     string
	Type     string
	Repeated bool
}

type MessageVisitor struct {
	proto.NoopVisitor
}

func (m MessageVisitor) VisitNormalField(i *proto.NormalField) {
	ProtoField.Name = i.Field.Name
	ProtoField.Type = i.Field.Type
	ProtoField.Repeated = i.Repeated
}

func GenCommentString(comments []string, space bool) string {
	var commentsString strings.Builder
	var spaceString string
	if space {
		spaceString = "  "
	}

	for _, v := range comments {
		commentsString.WriteString(fmt.Sprintf("%s// %s\n", spaceString, v))
	}
	return commentsString.String()
}
