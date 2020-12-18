package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type syntaxState struct {
	syntax string
}

func newSyntaxState(syntax string) *syntaxState {
	return &syntaxState{
		syntax: syntax,
	}
}

func (s *syntaxState) process() (string, error) {
	s.syntax = strings.ReplaceAll(s.syntax, `"`, "")
	colonIndex := strings.Index(s.syntax, ":")
	if colonIndex <= 0 {
		return "", errors.New("expected ':' near syntax")
	}

	flagVIndex := strings.Index(s.syntax, "v")
	if flagVIndex < 0 {
		return "", errors.New("expected version spec after ':'")
	}

	version := s.syntax[flagVIndex:]
	versionNumber := strings.TrimPrefix(version, "v")
	_, err := strconv.ParseUint(versionNumber, 10, 64)
	if err != nil {
		return "", fmt.Errorf("expected version spec after ':', but found %s", version)
	}

	return version, nil
}

func (s *syntaxState) writeInfo(api *spec.ApiSpec, attrs map[string]string) error {
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
