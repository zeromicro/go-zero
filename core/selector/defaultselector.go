package selector

import (
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/balancer"
)

const DefaultSelector = "defaultSelector"

var (
	_                           Selector = (*defaultSelector)(nil)
	colorAttributeKey                    = attribute.Key("selector.color")
	candidateColorsAttributeKey          = attribute.Key("selector.clientColors")
)

func init() {
	Register(defaultSelector{})
}

// defaultSelector a default selector.
type defaultSelector struct{}

func (d defaultSelector) Name() string {
	return DefaultSelector
}

func (d defaultSelector) Select(conns []Conn, info balancer.PickInfo) []Conn {
	clientColorsVal := ColorsFromContext(info.Ctx)
	clientColors := clientColorsVal.Colors()
	spanCtx := trace.SpanFromContext(info.Ctx)
	spanCtx.SetAttributes(candidateColorsAttributeKey.StringSlice(clientColors))

	connMap := d.genColor2ConnsMap(conns)
	if len(connMap) == 0 {
		// There is no dyed connection on the server.
		logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("clientColors", clientColorsVal.String()))
		return conns
	}

	newConns := make([]Conn, 0, len(conns))
	for _, clientColor := range clientColors {
		if v, yes := connMap[clientColor]; yes {
			newConns = append(newConns, v...)
		}

		if len(newConns) != 0 {
			spanCtx.SetAttributes(colorAttributeKey.String(clientColor))
			logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("color", clientColor), logx.Field("clientColors", clientColorsVal.String()))
			break
		}
	}

	return newConns
}

func (d defaultSelector) genColor2ConnsMap(conns []Conn) map[string][]Conn {
	m := map[string][]Conn{}
	for _, conn := range conns {
		address := conn.Address()
		serverColorsVal := address.BalancerAttributes.Value("colors")

		switch v := serverColorsVal.(type) {
		case *Colors:
			if v != nil {
				for _, color := range v.Colors() {
					m[color] = append(m[color], conn)
				}
			}
		}

	}

	return m
}
