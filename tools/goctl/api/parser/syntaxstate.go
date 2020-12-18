package parser

import (
	"fmt"
	"strconv"
	"strings"
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
	fields := strings.Fields(s.syntax)
	if len(fields) != 3 {
		return "", fmt.Errorf("expected syntax, but found %s", s.syntax)
	}

	if fields[0] != tokenSyntax {
		return "", fmt.Errorf("expected syntax, but found %s", fields[0])
	}

	if fields[1] != "=" {
		return "", fmt.Errorf("expected '=' after syntax, but found %s", fields[1])
	}

	version := fields[2]
	version = strings.TrimPrefix(version, `"`)
	version = strings.TrimSuffix(version, `"`)
	if !strings.HasPrefix(version, "v") {
		return "", fmt.Errorf("expected version after '=', but found %s", version)
	}

	versionNumber := strings.TrimPrefix(version, "v")
	_, err := strconv.ParseUint(versionNumber, 10, 64)
	if err != nil {
		return "", fmt.Errorf("expected version after '=', but found %s", version)
	}

	return version, nil
}
