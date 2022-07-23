package md

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCarrier struct {
	err error
	md  Metadata
}

func (c *mockCarrier) Carrier() (Metadata, error) {
	return c.md, c.err
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
	assert.EqualValues(t, []string{"a", "b", "c"}, metadata.Keys())
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
	metadata = Metadata{"a": {"a1"}, "b": {}}
	assert.Equal(t, []string{"a1"}, metadata.Values("a"))
	assert.Equal(t, []string{"a1"}, metadata.Values("A"))
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

func TestMetadata_Carrier(t *testing.T) {
	metadata := Metadata{}
	m, err := metadata.Carrier()
	assert.NoError(t, err)
	assert.Equal(t, metadata, m)

	metadata = Metadata{"a": {"a1"}, "b": {}}
	m, err = metadata.Carrier()
	assert.NoError(t, err)
	assert.Equal(t, metadata, m)
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
	ctx := NewContext(context.Background(), &mockCarrier{md: Metadata{"a": {"a1"}}})
	assert.Equal(t, Metadata{"a": {"a1"}}, ctx.Value(mdKey{}))

	ctx = NewContext(context.Background(), &mockCarrier{md: Metadata{"a": {"a1"}}, err: errors.New("any")})
	assert.Equal(t, nil, ctx.Value(mdKey{}))
}

func TestValuesFromContext(t *testing.T) {
	md := Metadata{"a": {"a1"}}
	assert.Equal(t, []string{"a1"}, ValuesFromContext(context.WithValue(context.Background(), mdKey{}, md), "a"))
	assert.Equal(t, []string{"a1"}, ValuesFromContext(context.WithValue(context.Background(), mdKey{}, md), "A"))
	assert.Equal(t, []string(nil), ValuesFromContext(context.WithValue(context.Background(), mdKey{}, md), "b"))
}
