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

func (d defaultSelector) Select(conns []Conn, info balancer.PickInfo) []Conn {
	clientColorsVal := ColorsFromContext(info.Ctx)
	if clientColorsVal.Empty() {
		return d.getNoColorsConns(conns)
	}

	clientColors := clientColorsVal.Colors()
	spanCtx := trace.SpanFromContext(info.Ctx)
	spanCtx.SetAttributes(candidateColorsAttributeKey.StringSlice(clientColors))

	var newConns []Conn
	m := d.genColor2ConnsMap(conns)
	for _, clientColor := range clientColors {
		if v, yes := m[clientColor]; yes {
			newConns = append(newConns, v...)
		}

		if len(newConns) != 0 {
			spanCtx.SetAttributes(colorAttributeKey.String(clientColor))
			logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("color", clientColor), logx.Field("clientColors", clientColorsVal.String()))
			break
		}
	}

	if len(newConns) == 0 {
		logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("clientColors", clientColorsVal.String()))
	}

	return newConns
}

func (d defaultSelector) Name() string {
	return DefaultSelector
}

func (d defaultSelector) genColor2ConnsMap(conns []Conn) map[string][]Conn {
	m := map[string][]Conn{}
	for _, conn := range conns {
		address := conn.Address()
		serverColorsVal := address.BalancerAttributes.Value("colors")
		if serverColorsVal == nil {
			continue
		}

		var serverColors []string
		if c, ok := serverColorsVal.(*Colors); ok && c != nil {
			serverColors = c.Colors()
		}

		for _, color := range serverColors {
			m[color] = append(m[color], conn)
		}
	}

	return m
}

func (d defaultSelector) getNoColorsConns(conns []Conn) []Conn {
	var newConns []Conn
	for _, conn := range conns {
		address := conn.Address()
		colorsVal := address.BalancerAttributes.Value("colors")
		if colorsVal == nil {
			continue
		}
		c, ok := colorsVal.(*Colors)
		if !ok {
			continue
		}

		if c != nil || c.Size() == 0 {
			newConns = append(newConns, conn)
		}
	}

	return newConns
}
