package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

type serviceState struct {
	*baseState
	annos []spec.Annotation
}

func newServiceState(state *baseState, annos []spec.Annotation) state {
	return &serviceState{
		baseState: state,
		annos:     annos,
	}
}

func (s *serviceState) process(api *spec.ApiSpec) (state, error) {
	var name string
	var routes []spec.Route
	parser := &serviceEntityParser{
		acceptName: func(n string) {
			name = n
		},
		acceptRoute: func(route spec.Route) {
			routes = append(routes, route)
		},
	}
	ent := newEntity(s.baseState, api, parser)
	if err := ent.process(); err != nil {
		return nil, err
	}

	api.Service = spec.Service{
		Name:        name,
		Annotations: append(api.Service.Annotations, s.annos...),
		Routes:      append(api.Service.Routes, routes...),
		Groups: append(api.Service.Groups, spec.Group{
			Annotations: s.annos,
			Routes:      routes,
		}),
	}

	return newRootState(s.r, s.lineNumber), nil
}

type serviceEntityParser struct {
	acceptName  func(name string)
	acceptRoute func(route spec.Route)
}

func (p *serviceEntityParser) parseLine(line string, api *spec.ApiSpec, annos []spec.Annotation) error {
	line = strings.TrimSpace(line)

	var buffer = new(bytes.Buffer)
	buffer.WriteString(line)
	reader := bufio.NewReader(buffer)
	var builder strings.Builder
	var fields []string
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch {
		case isSpace(ch), ch == leftParenthesis, ch == rightParenthesis, ch == semicolon:
			if builder.Len() == 0 {
				continue
			}
			token := builder.String()
			builder.Reset()
			fields = append(fields, token)
		default:
			builder.WriteRune(ch)
		}
	}

	if len(fields) < 3 {
		return fmt.Errorf("wrong line %q, %q", line, routeSyntax)
	}

	method := fields[0]
	path := fields[1]
	req := fields[2]
	var returns string
	if len(fields) > 4 {
		returns = fields[4]
		if fields[3] != returnsTag {
			return fmt.Errorf("wrong line %q, %q", line, routeSyntax)
		}
	}

	p.acceptRoute(spec.Route{
		Annotations:  annos,
		Method:       method,
		Path:         path,
		RequestType:  GetType(api, req),
		ResponseType: GetType(api, returns),
	})

	return nil
}

func (p *serviceEntityParser) setEntityName(name string) {
	p.acceptName(name)
}
