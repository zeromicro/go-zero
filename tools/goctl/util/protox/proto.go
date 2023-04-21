// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Optional bool
	Sequence int
}

type MessageVisitor struct {
	proto.NoopVisitor
}

func (m MessageVisitor) VisitNormalField(i *proto.NormalField) {
	ProtoField.Name = i.Field.Name
	ProtoField.Type = i.Field.Type
	ProtoField.Repeated = i.Repeated
	ProtoField.Optional = i.Optional
	ProtoField.Sequence = i.Sequence
}

func (m MessageVisitor) VisitMapField(i *proto.MapField) {
	ProtoField.Name = i.Field.Name
	ProtoField.Type = fmt.Sprintf("map<%s,%s>", i.KeyType, i.Field.Type)
	ProtoField.Sequence = i.Sequence
	ProtoField.Repeated = false
	ProtoField.Optional = false
}

func (m MessageVisitor) VisitEnumField(i *proto.EnumField) {
	ProtoField.Name = i.Name
	ProtoField.Type = ""
	ProtoField.Sequence = i.Integer
	ProtoField.Repeated = false
	ProtoField.Optional = false
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
