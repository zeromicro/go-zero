package trace

import (
	"context"

	"github.com/zeromicro/go-zero/core/logc"
	"go.opentelemetry.io/otel/baggage"
)

// GetBaggageValue get baggage info from context, if key not exists, return "", false.
func GetBaggageValue(ctx context.Context, key string) (string, bool) {
	b := baggage.FromContext(ctx)
	m := b.Member(key)

	if m.Value() == "" {
		return "", false
	}

	return m.Value(), true
}

// WithBaggage append baggage by string key val.
func WithBaggage(parent context.Context, key, val string) context.Context {
	member, err := baggage.NewMember(key, val)
	if err != nil {
		logc.Error(parent, err)
		return parent
	}

	b := baggage.FromContext(parent)
	b, err = b.SetMember(member)
	if err != nil {
		logc.Error(parent, err)
		return parent
	}

	return baggage.ContextWithBaggage(parent, b)
}

// AddBaggagesFromMap append map kvs to current ctx baggage.
func AddBaggagesFromMap(parent context.Context, mp map[string]string) context.Context {
	b := baggage.FromContext(parent)

	for k, v := range mp {
		m, err := baggage.NewMember(k, v)
		if err != nil {
			logc.Error(parent, err)
			return parent
		}

		b, err = b.SetMember(m)
		if err != nil {
			logc.Error(parent, err)
			return parent
		}
	}

	return baggage.ContextWithBaggage(parent, b)
}
