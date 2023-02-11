package kube

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/zrpc/resolver/internal/targets"
	"google.golang.org/grpc/resolver"
)

const (
	colon            = ":"
	defaultNamespace = "default"
	queryNonBlock    = "nonBlock"
)

var emptyService Service

// Service represents a service with namespace, name and port.
type Service struct {
	Namespace string
	Name      string
	Port      int
	NonBlock  bool
}

// ParseTarget parses the resolver.Target.
func ParseTarget(target resolver.Target) (Service, error) {
	var service Service
	service.Namespace = targets.GetAuthority(target)
	if len(service.Namespace) == 0 {
		service.Namespace = defaultNamespace
	}

	endpoints := targets.GetEndpoints(target)
	if strings.Contains(endpoints, colon) {
		segs := strings.SplitN(endpoints, colon, 2)
		if len(segs) < 2 {
			return emptyService, fmt.Errorf("bad endpoint: %s", endpoints)
		}

		service.Name = segs[0]
		port, err := strconv.Atoi(segs[1])
		if err != nil {
			return emptyService, err
		}

		service.Port = port
	} else {
		service.Name = endpoints
	}

	return parseNonBlock(target, service)
}

// parseNonBlock
func parseNonBlock(target resolver.Target, service Service) (Service, error) {

	values, err := url.ParseQuery(target.URL.RawQuery)
	if err != nil {
		return emptyService, err
	}

	if vales, ok := values[queryNonBlock]; ok {
		nonBlock, err := strconv.ParseBool(vales[0])
		if err != nil {
			return emptyService, err
		}
		service.NonBlock = nonBlock
	}

	return service, nil
}
