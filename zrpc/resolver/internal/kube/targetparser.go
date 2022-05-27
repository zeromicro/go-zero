package kube

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/zrpc/resolver/internal/targets"
	"google.golang.org/grpc/resolver"
)

const (
	colon            = ":"
	defaultNamespace = "default"
)

var emptyService Service

// Service represents a service with namespace, name and port.
type Service struct {
	Namespace string
	Name      string
	Port      int
}

// ParseTarget parses the resolver.Target.
func ParseTarget(target resolver.Target) (Service, error) {
	var service Service
	service.Namespace = targets.GetAuthority(target)
	if len(service.Namespace) == 0 {
		service.Namespace = defaultNamespace
	}

	endpoints := targets.GetEndpoints(target)
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

	return service, nil
}
