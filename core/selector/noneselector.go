package selector

import "google.golang.org/grpc/balancer"

type noneSelector struct{}

func (n noneSelector) Name() string {
	return ""
}

func (n noneSelector) Select(conns []Conn, _ balancer.PickInfo) []Conn { return conns }
