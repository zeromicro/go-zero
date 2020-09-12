package parser

import (
	"fmt"
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
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("wrong line %q", line)
	}

	method := fields[0]
	pathAndRequest := fields[1]
	pos := strings.Index(pathAndRequest, "(")
	if pos < 0 {
		return fmt.Errorf("wrong line %q", line)
	}
	path := strings.TrimSpace(pathAndRequest[:pos])
	pathAndRequest = pathAndRequest[pos+1:]
	pos = strings.Index(pathAndRequest, ")")
	if pos < 0 {
		return fmt.Errorf("wrong line %q", line)
	}
	req := pathAndRequest[:pos]
	var returns string
	if len(fields) > 2 {
		returns = fields[2]
	}
	returns = strings.ReplaceAll(returns, "returns", "")
	returns = strings.ReplaceAll(returns, "(", "")
	returns = strings.ReplaceAll(returns, ")", "")
	returns = strings.TrimSpace(returns)

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
