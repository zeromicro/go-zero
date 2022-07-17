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
	candidateColorsAttributeKey          = attribute.Key("selector.candidateColors")
)

func init() {
	Register(defaultSelector{})
}

type defaultSelector struct{}

func (d defaultSelector) Select(conns []Conn, info balancer.PickInfo) []Conn {
	clientColorsVal, ok := ColorsFromContext(info.Ctx)
	if !ok {
		return d.getNoColorConns(conns)
	}
	clientColors := clientColorsVal.Colors()

	spanCtx := trace.SpanFromContext(info.Ctx)
	spanCtx.SetAttributes(candidateColorsAttributeKey.StringSlice(clientColors))

	newConns := make([]Conn, 0, len(conns))
	for _, clientColor := range clientColors {
		for _, conn := range conns {
			address := conn.Address()
			serverColorsVal := address.BalancerAttributes.Value("colors")
			if serverColorsVal == nil {
				continue
			}
			c, b := serverColorsVal.(*Colors)
			if !b || c.Size() == 0 {
				continue
			}

			serverColors := c.Colors()
			for _, serverColor := range serverColors {
				if clientColor == serverColor {
					newConns = append(newConns, conn)
				}
			}

		}

		if len(newConns) != 0 {
			spanCtx.SetAttributes(colorAttributeKey.String(clientColor))
			logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("color", clientColor), logx.Field("candidateColors", clientColorsVal.String()))

			break
		}
	}

	if len(newConns) == 0 {
		logx.WithContext(info.Ctx).Infow("flow dyeing", logx.Field("candidateColors", clientColorsVal.String()))
	}

	return newConns
}

func (d defaultSelector) Name() string {
	return DefaultSelector
}

func (d defaultSelector) getNoColorConns(conns []Conn) []Conn {
	var newConns []Conn
	for _, conn := range conns {
		address := conn.Address()
		colorsVal := address.BalancerAttributes.Value("colors")
		if colorsVal == nil {
			newConns = append(newConns, conn)
		}
		c, b := colorsVal.(*Colors)
		if !b || c.Size() == 0 {
			newConns = append(newConns, conn)
		}
	}

	return newConns
}
