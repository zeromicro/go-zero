package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
)

func (p *Parser) validate(api *spec.ApiSpec) (err error) {
	var builder strings.Builder
	for _, tp := range api.Types {
		if ok, name := p.validateDuplicateProperty(tp); !ok {
			fmt.Fprintf(&builder, `duplicate property "%s" of type "%s"`+"\n", name, tp.Name)
		}
	}
	if ok, info := p.validateDuplicateRouteHandler(api); !ok {
		fmt.Fprintf(&builder, info)
	}
	if len(builder.String()) > 0 {
		return errors.New(builder.String())
	}
	return nil
}

func (p *Parser) validateDuplicateProperty(tp spec.Type) (bool, string) {
	var names []string
	for _, member := range tp.Members {
		if stringx.Contains(names, member.Name) {
			return false, member.Name
		} else {
			names = append(names, member.Name)
		}
	}
	return true, ""
}

func (p *Parser) validateDuplicateRouteHandler(api *spec.ApiSpec) (bool, string) {
	var names []string
	for _, r := range api.Service.Routes() {
		handler, ok := util.GetAnnotationValue(r.Annotations, "server", "handler")
		if !ok {
			return false, fmt.Sprintf("missing handler annotation for %s", r.Path)
		}
		if stringx.Contains(names, handler) {
			return false, fmt.Sprintf(`duplicated handler for name "%s"`, handler)
		} else {
			names = append(names, handler)
		}
	}
	return true, ""
}
