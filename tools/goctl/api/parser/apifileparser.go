package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	tokenInfo              = "info"
	tokenImport            = "import"
	tokenType              = "type"
	tokenService           = "service"
	tokenServiceAnnotation = "@server"
)

type (
	ApiStruct struct {
		Info             string
		Type             string
		Service          string
		Imports          string
		serviceBeginLine int
	}

	apiFileState interface {
		process(api *ApiStruct, token string) (apiFileState, error)
	}

	apiRootState struct {
		*baseState
	}

	apiInfoState struct {
		*baseState
	}

	apiImportState struct {
		*baseState
	}

	apiTypeState struct {
		*baseState
	}

	apiServiceState struct {
		*baseState
	}
)

func ParseApi(src string) (*ApiStruct, error) {
	var buffer = new(bytes.Buffer)
	buffer.WriteString(src)
	api := new(ApiStruct)
	var lineNumber = api.serviceBeginLine
	apiFile := baseState{r: bufio.NewReader(buffer), lineNumber: &lineNumber}
	st := apiRootState{&apiFile}
	for {
		st, err := st.process(api, "")
		if err == io.EOF {
			return api, nil
		}
		if err != nil {
			return nil, fmt.Errorf("near line: %d, %s", lineNumber, err.Error())
		}
		if st == nil {
			return api, nil
		}
	}
}

func (s *apiRootState) process(api *ApiStruct, token string) (apiFileState, error) {
	var builder strings.Builder
	for {
		ch, err := s.readSkipComment()
		if err != nil {
			return nil, err
		}

		switch {
		case isSpace(ch) || isNewline(ch) || ch == leftParenthesis:
			token := builder.String()
			token = strings.TrimSpace(token)
			if len(token) == 0 {
				continue
			}

			builder.Reset()
			switch token {
			case tokenInfo:
				info := apiInfoState{s.baseState}
				return info.process(api, token+string(ch))
			case tokenImport:
				tp := apiImportState{s.baseState}
				return tp.process(api, token+string(ch))
			case tokenType:
				ty := apiTypeState{s.baseState}
				return ty.process(api, token+string(ch))
			case tokenService:
				server := apiServiceState{s.baseState}
				return server.process(api, token+string(ch))
			case tokenServiceAnnotation:
				server := apiServiceState{s.baseState}
				return server.process(api, token+string(ch))
			default:
				if strings.HasPrefix(token, "//") {
					continue
				}
				return nil, errors.New(fmt.Sprintf("invalid token %s at line %d", token, *s.lineNumber))
			}
		default:
			builder.WriteRune(ch)
		}
	}
}

func (s *apiInfoState) process(api *ApiStruct, token string) (apiFileState, error) {
	for {
		line, err := s.readLine()
		if err != nil {
			return nil, err
		}

		api.Info += "\n" + token + line
		token = ""
		if strings.TrimSpace(line) == string(rightParenthesis) {
			return &apiRootState{s.baseState}, nil
		}
	}
}

func (s *apiImportState) process(api *ApiStruct, token string) (apiFileState, error) {
	line, err := s.readLine()
	if err != nil {
		return nil, err
	}

	line = token + line
	if len(strings.Fields(line)) != 2 {
		return nil, errors.New("import syntax error: " + line)
	}

	api.Imports += "\n" + line
	return &apiRootState{s.baseState}, nil
}

func (s *apiTypeState) process(api *ApiStruct, token string) (apiFileState, error) {
	var blockCount = 0
	for {
		line, err := s.readLine()
		if err != nil {
			return nil, err
		}

		api.Type += "\n\n" + token + line
		token = ""
		line = strings.TrimSpace(line)
		line = removeComment(line)
		if strings.HasSuffix(line, leftBrace) {
			blockCount++
		}
		if strings.HasSuffix(line, string(leftParenthesis)) {
			blockCount++
		}
		if strings.HasSuffix(line, string(rightBrace)) {
			blockCount--
		}
		if strings.HasSuffix(line, string(rightParenthesis)) {
			blockCount--
		}

		if blockCount == 0 {
			return &apiRootState{s.baseState}, nil
		}
	}
}

func (s *apiServiceState) process(api *ApiStruct, token string) (apiFileState, error) {
	var blockCount = 0
	for {
		line, err := s.readLineSkipComment()
		if err != nil {
			return nil, err
		}

		line = token + line
		token = ""
		api.Service += "\n" + line
		line = strings.TrimSpace(line)
		line = removeComment(line)
		if strings.HasSuffix(line, leftBrace) {
			blockCount++
		}
		if strings.HasSuffix(line, string(leftParenthesis)) {
			blockCount++
		}
		if line == string(rightBrace) {
			blockCount--
		}
		if line == string(rightParenthesis) {
			blockCount--
		}

		if blockCount == 0 {
			return &apiRootState{s.baseState}, nil
		}
	}
}

func removeComment(line string) string {
	var commentIdx = strings.Index(line, "//")
	if commentIdx >= 0 {
		return line[:commentIdx]
	}
	return line
}
