package md

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCarrier struct {
	md  map[string][]string
	err error
}

func (m *mockCarrier) Extract(ctx context.Context) (context.Context, error) {
	if m.err != nil {
		return ctx, m.err
	}

	metadata := FromContext(ctx)
	metadata = metadata.Clone()
	for k, v := range m.md {
		metadata.Append(strings.ToLower(k), v...)
	}

	return NewContext(ctx, metadata), m.err
}

func (m *mockCarrier) Inject(ctx context.Context) error {
	if m.err != nil {
		return m.err
	}

	metadata := FromContext(ctx)
	for k, v := range metadata {
		m.md[strings.ToLower(k)] = v
	}

	return nil
}

func TestMetadata_Append(t *testing.T) {
	metadata := Metadata{}
	metadata.Append("a", "a1")
	assert.EqualValues(t, map[string][]string{"a": {"a1"}}, metadata)
	metadata.Append("a", "a2")
	assert.EqualValues(t, map[string][]string{"a": {"a1", "a2"}}, metadata)
	metadata.Append("b", "b1", "b2")
	assert.EqualValues(t, map[string][]string{"a": {"a1", "a2"}, "b": {"b1", "b2"}}, metadata)
}

func TestMetadata_Keys(t *testing.T) {
	metadata := Metadata{}
	assert.Equal(t, []string(nil), metadata.Keys())
	metadata = Metadata{"a": {}, "b": {}, "c": {"1"}}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, metadata.Keys())
}

func TestMetadata_Set(t *testing.T) {
	metadata := Metadata{}
	assert.Len(t, metadata, 0)
	metadata.Set("a", "a1")
	assert.EqualValues(t, map[string][]string{"a": {"a1"}}, metadata)
	metadata.Set("a", "a1", "a2")
	assert.EqualValues(t, map[string][]string{"a": {"a1", "a2"}}, metadata)
	metadata.Set("A", "a1")
	assert.EqualValues(t, map[string][]string{"a": {"a1"}}, metadata)
}

func TestMetadata_Get(t *testing.T) {
	metadata := Metadata{}
	assert.Len(t, metadata, 0)
	metadata = Metadata{"a": {"a1"}, "b": {}, "c": {"c1", "c2"}}
	assert.Equal(t, "a1", metadata.Get("a"))
	assert.Equal(t, "a1", metadata.Get("A"))
	assert.Equal(t, "", metadata.Get("b"))
	assert.Equal(t, "", metadata.Get("B"))
	assert.Equal(t, "c1", metadata.Get("c"))
	assert.Equal(t, "c1", metadata.Get("C"))
	assert.Equal(t, "", metadata.Get("D"))
}

func TestMetadata_Values(t *testing.T) {
	metadata := Metadata{}
	assert.Len(t, metadata, 0)
	metadata = Metadata{"a": {"a1", "a2"}, "b": {}}
	assert.Equal(t, []string{"a1", "a2"}, metadata.Values("a"))
	assert.Equal(t, []string{"a1", "a2"}, metadata.Values("A"))
	assert.Equal(t, []string{}, metadata.Values("b"))
	assert.Equal(t, []string{}, metadata.Values("B"))
	assert.Equal(t, []string(nil), metadata.Values("C"))
}

func TestMetadata_Delete(t *testing.T) {
	metadata := Metadata{"a": {"a1"}, "b": {}}
	assert.EqualValues(t, map[string][]string{"a": {"a1"}, "b": {}}, metadata)
	metadata.Delete("a")
	assert.EqualValues(t, map[string][]string{"b": {}}, metadata)
	metadata.Delete("B")
	assert.EqualValues(t, map[string][]string{}, metadata)
}

func TestMetadata_Clone(t *testing.T) {
	metadata := Metadata{}
	assert.Equal(t, metadata, metadata.Clone())

	metadata = Metadata{"a": {"a1"}, "b": {}}
	assert.Equal(t, metadata, metadata.Clone())
}

func TestFromContext(t *testing.T) {
	assert.Equal(t, Metadata{}, FromContext(context.Background()))
	assert.Equal(t, Metadata{"a": {"a1"}}, FromContext(context.WithValue(context.Background(), mdKey{}, Metadata{"a": {"a1"}})))
}

func TestNewContext(t *testing.T) {
	ctx := NewContext(context.Background(), Metadata{"a": {"a1"}})
	assert.Equal(t, Metadata{"a": {"a1"}}, ctx.Value(mdKey{}))

	ctx = NewContext(context.Background(), Metadata{})
	assert.Equal(t, Metadata{}, ctx.Value(mdKey{}))
}

func TestValuesFromContext(t *testing.T) {
	md := Metadata{"a": {"a1"}}
	assert.Equal(t, []string{"a1"}, ValuesFromContext(context.WithValue(context.Background(), mdKey{}, md), "a"))
	assert.Equal(t, []string{"a1"}, ValuesFromContext(context.WithValue(context.Background(), mdKey{}, md), "A"))
	assert.Equal(t, []string(nil), ValuesFromContext(context.WithValue(context.Background(), mdKey{}, md), "b"))
}

func TestExtract(t *testing.T) {
	t.Run("no err", func(t *testing.T) {
		m := &mockCarrier{md: map[string][]string{"a": {"a1"}}}
		ctx := Extract(context.Background(), m)
		assert.EqualValues(t, map[string][]string{"a": {"a1"}}, FromContext(ctx))
	})

	t.Run("err", func(t *testing.T) {
		m := &mockCarrier{md: map[string][]string{"a": {"a1"}}, err: errors.New("any")}
		ctx := Extract(context.Background(), m)
		assert.EqualValues(t, map[string][]string{}, FromContext(ctx))
	})
}

func TestInjection(t *testing.T) {
	t.Run("no err", func(t *testing.T) {
		m := map[string][]string{}
		Inject(NewContext(context.Background(), Metadata{"a": {"a1"}}), &mockCarrier{md: m})
		assert.Equal(t, map[string][]string{"a": {"a1"}}, m)
	})

	t.Run("no err", func(t *testing.T) {
		m := map[string][]string{}
		Inject(NewContext(context.Background(), Metadata{"a": {"a1"}}), &mockCarrier{md: m, err: errors.New("any")})
		assert.Equal(t, map[string][]string{}, m)
	})
}
