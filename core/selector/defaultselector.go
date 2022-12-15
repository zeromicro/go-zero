package selector

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/md"
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

	clientColors := md.ValuesFromContext(info.Ctx, "colors")
	spanCtx := trace.SpanFromContext(info.Ctx)
	trace.SpanFromContext(info.Ctx).SetAttributes(candidateColorsAttributeKey.StringSlice(clientColors))

	connMap := d.genColor2ConnsMap(conns)
	if len(connMap) == 0 {
		// There is no dyed connection on the server.
		logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("clientColors", clientColors))
		return conns
	}

	newConns := make([]Conn, 0, len(conns))
	for _, clientColor := range clientColors {
		if v, yes := connMap[clientColor]; yes {
			newConns = append(newConns, v...)
		}

		if len(newConns) != 0 {
			spanCtx.SetAttributes(colorAttributeKey.String(clientColor))
			logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("color", clientColor), logx.Field("clientColors", clientColors))
			break
		}
	}

	return newConns
}

func (d defaultSelector) genColor2ConnsMap(conns []Conn) map[string][]Conn {
	m := map[string][]Conn{}
	for _, conn := range conns {
		address := conn.Address()
		metadataVal := address.BalancerAttributes.Value("metadata")

		switch metadata := metadataVal.(type) {
		case md.Metadata:
			colors := metadata.Values("colors")
			for _, color := range colors {
				m[color] = append(m[color], conn)
			}
		}
	}

	return m
}
