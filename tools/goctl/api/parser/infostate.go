package parser

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const (
	titleTag   = "title"
	descTag    = "desc"
	versionTag = "version"
	authorTag  = "author"
	emailTag   = "email"
)

type infoState struct {
	*baseState
	innerState int
}

func newInfoState(st *baseState) state {
	return &infoState{
		baseState:  st,
		innerState: startState,
	}
}

func (s *infoState) process(api *spec.ApiSpec) (state, error) {
	attrs, err := s.parseProperties()
	if err != nil {
		return nil, err
	}

	if err := s.writeInfo(api, attrs); err != nil {
		return nil, err
	}

	return newRootState(s.r, s.lineNumber), nil
}

func (s *infoState) writeInfo(api *spec.ApiSpec, attrs map[string]string) error {
	for k, v := range attrs {
		switch k {
		case titleTag:
			api.Info.Title = strings.TrimSpace(v)
		case descTag:
			api.Info.Desc = strings.TrimSpace(v)
		case versionTag:
			api.Info.Version = strings.TrimSpace(v)
		case authorTag:
			api.Info.Author = strings.TrimSpace(v)
		case emailTag:
			api.Info.Email = strings.TrimSpace(v)
		default:
			return fmt.Errorf("unknown directive %q in %q section", k, infoDirective)
		}
	}

	return nil
}
