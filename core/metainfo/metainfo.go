package metainfo

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/propagation"

	"github.com/zeromicro/go-zero/core/collection"
)

const (
	// PrefixPass means that header/metadata key with this prefix will be passed to the other servers.
	PrefixPass = "x-pass-"

	lenPP = len(PrefixPass)
)

var (
	// CustomKeysMapPropagator impl propagation.TextMapPropagator for custom keys passing.
	CustomKeysMapPropagator propagation.TextMapPropagator = (*customKeysPropagator)(nil)

	ctxKey         ctxKeyType
	customKeyStore = contextKeyStore{
		keyArr: make([]string, 0),
		keySet: collection.NewSet(),
	}
)

type (
	customKeysPropagator struct{}

	ctxKeyType      struct{}
	contextKeyStore struct {
		keyArr []string
		keySet *collection.Set
	}
)

// RegisterCustomKeys register custom keys globally.
// Key must be lowercase.
// Should only be called once before application start.
func RegisterCustomKeys(keys []string) {
	for _, k := range keys {
		kk := strings.ToLower(k)
		if k != kk {
			panic("custom key only support lowercase")
		}
		customKeyStore.keySet.AddStr(k)
	}
	customKeyStore.keyArr = customKeyStore.keySet.KeysStr()
}

// for test only
func reset() {
	customKeyStore = contextKeyStore{
		keyArr: make([]string, 0),
		keySet: collection.NewSet(),
	}
}

func getMap(ctx context.Context) map[string]string {
	if ctx != nil {
		if val, ok := ctx.Value(ctxKey).(map[string]string); ok {
			return val
		}
	}

	return make(map[string]string, 0)
}

func setMap(ctx context.Context, m map[string]string) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, ctxKey, m)
}

// GetMapFromContext retrieves all custom keys and values from context.
func GetMapFromContext(ctx context.Context) map[string]string {
	mp := getMap(ctx)

	if len(mp) > 0 {
		m := make(map[string]string, len(mp))
		for k, v := range mp {
			m[k] = v
		}
		return m
	}

	return mp
}

// GetMapFromPropagator retrieves all custom keys and values from propagation carrier.
func GetMapFromPropagator(carrier propagation.TextMapCarrier) map[string]string {
	mp := make(map[string]string)
	for _, k := range carrier.Keys() {
		kk := strings.ToLower(k)

		if customKeyStore.keySet.Contains(kk) || (len(kk) > lenPP && strings.HasPrefix(kk, PrefixPass)) {
			v := carrier.Get(kk)
			if len(v) > 0 {
				mp[kk] = v
			}
		}
	}
	return mp
}

// Inject impl TextMapPropagator for customKeysPropagator.
func (c *customKeysPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	mp := getMap(ctx)
	for k, v := range mp {
		if len(v) > 0 {
			carrier.Set(k, v)
		}
	}
}

// Extract impl TextMapPropagator for customKeysPropagator.
func (c *customKeysPropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	mp := getMap(ctx)

	cmp := GetMapFromPropagator(carrier)
	if len(cmp) == 0 {
		return ctx
	}

	for k, v := range cmp {
		mp[k] = v
	}

	return setMap(ctx, mp)
}

// Fields not used
func (c *customKeysPropagator) Fields() []string {
	return nil
}
