package parser

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

type typeState struct {
	*baseState
	annos []spec.Annotation
}

func newTypeState(state *baseState, annos []spec.Annotation) state {
	return &typeState{
		baseState: state,
		annos:     annos,
	}
}

func (s *typeState) process(api *spec.ApiSpec) (state, error) {
	var name string
	var members []spec.Member
	parser := &typeEntityParser{
		acceptName: func(n string) {
			name = n
		},
		acceptMember: func(member spec.Member) {
			members = append(members, member)
		},
	}
	ent := newEntity(s.baseState, api, parser)
	if err := ent.process(); err != nil {
		return nil, err
	}

	api.Types = append(api.Types, spec.Type{
		Name:        name,
		Annotations: s.annos,
		Members:     members,
	})

	return newRootState(s.r, s.lineNumber), nil
}

type typeEntityParser struct {
	acceptName   func(name string)
	acceptMember func(member spec.Member)
}

func (p *typeEntityParser) parseLine(line string, api *spec.ApiSpec, annos []spec.Annotation) error {
	index := strings.Index(line, "//")
	comment := ""
	if index >= 0 {
		comment = line[index+2:]
		line = strings.TrimSpace(line[:index])
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	if len(fields) == 1 {
		p.acceptMember(spec.Member{
			Annotations: annos,
			Name:        fields[0],
			Type:        fields[0],
			IsInline:    true,
		})
		return nil
	}
	name := fields[0]
	tp := fields[1]
	var tag string
	if len(fields) > 2 {
		tag = fields[2]
	} else {
		tag = fmt.Sprintf("`json:\"%s\"`", util.Untitle(name))
	}

	p.acceptMember(spec.Member{
		Annotations: annos,
		Name:        name,
		Type:        tp,
		Tag:         tag,
		Comment:     comment,
		IsInline:    false,
	})
	return nil
}

func (p *typeEntityParser) setEntityName(name string) {
	p.acceptName(name)
}
