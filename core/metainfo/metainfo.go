package metainfo

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	"go.opentelemetry.io/otel/propagation"
)

const (
	// PrefixPass means that header/metadata key with this prefix will be passed to the other servers.
	PrefixPass = "x-pass-"
	lenPP      = len(PrefixPass)
)

var (
	// CustomKeysMapPropagator implements propagation.TextMapPropagator for custom keys passing.
	CustomKeysMapPropagator propagation.TextMapPropagator = (*customKeysPropagator)(nil)

	ctxKey         ctxKeyType
	customKeyStore = newContextKeyStore()
)

type (
	customKeysPropagator struct{}

	ctxKeyType      struct{}
	contextKeyStore struct {
		keyArr []string
		keySet *collection.Set
	}
)

func newContextKeyStore() contextKeyStore {
	return contextKeyStore{
		keyArr: make([]string, 0),
		keySet: collection.NewSet(),
	}
}

// RegisterCustomKeys registers custom keys globally.
// Key must be lowercase.
// Should only be called once before application start.
func RegisterCustomKeys(keys []string) {
	for _, k := range keys {
		lowerKey := strings.ToLower(k)
		if k != lowerKey {
			panic("custom keys must be lowercase")
		}
		customKeyStore.keySet.AddStr(lowerKey)
	}
	customKeyStore.keyArr = customKeyStore.keySet.KeysStr()
}

// for test only
func reset() {
	customKeyStore = newContextKeyStore()
}

func getMap(ctx context.Context) map[string]string {
	if val, ok := ctx.Value(ctxKey).(map[string]string); ok {
		return val
	}

	return make(map[string]string)
}

func setMap(ctx context.Context, m map[string]string) context.Context {
	return context.WithValue(ctx, ctxKey, m)
}

// GetMapFromContext retrieves all custom keys and values from the context.
func GetMapFromContext(ctx context.Context) map[string]string {
	mp := getMap(ctx)
	if len(mp) == 0 {
		return mp
	}

	m := make(map[string]string, len(mp))
	for k, v := range mp {
		m[k] = v
	}
	return m
}

// GetMapFromPropagator retrieves all custom keys and values from the propagation carrier.
func GetMapFromPropagator(carrier propagation.TextMapCarrier) map[string]string {
	mp := make(map[string]string)
	for _, k := range carrier.Keys() {
		lowerKey := strings.ToLower(k)
		if customKeyStore.keySet.Contains(lowerKey) || (len(lowerKey) > lenPP &&
			strings.HasPrefix(lowerKey, PrefixPass)) {
			v := carrier.Get(lowerKey)
			if len(v) > 0 {
				mp[lowerKey] = v
			}
		}
	}
	return mp
}

// Inject implements TextMapPropagator for customKeysPropagator.
func (c *customKeysPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	mp := getMap(ctx)
	for k, v := range mp {
		if len(v) > 0 {
			carrier.Set(k, v)
		}
	}
}

// Extract implements TextMapPropagator for customKeysPropagator.
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

// Fields returns nil as it's not used.
func (c *customKeysPropagator) Fields() []string {
	return nil
}
