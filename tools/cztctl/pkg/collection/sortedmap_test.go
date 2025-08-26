package sortedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SortedMap(t *testing.T) {
	sm := New()
	t.Run("SetExpression", func(t *testing.T) {
		_, _, err := sm.SetExpression("")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		_, _, err = sm.SetExpression("foo")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		_, _, err = sm.SetExpression("foo= ")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		_, _, err = sm.SetExpression(" foo=")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		_, _, err = sm.SetExpression("foo =")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		_, _, err = sm.SetExpression("=")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		_, _, err = sm.SetExpression("=bar")
		assert.ErrorIs(t, err, ErrInvalidKVExpression)
		key, value, err := sm.SetExpression("foo=bar")
		assert.Nil(t, err)
		assert.Equal(t, "foo", key)
		assert.Equal(t, "bar", value)
		key, value, err = sm.SetExpression("foo=")
		assert.Nil(t, err)
		assert.Equal(t, value, sm.GetOr(key, ""))
		sm.Reset()
	})

	t.Run("SetKV", func(t *testing.T) {
		sm.SetKV("foo", "bar")
		assert.Equal(t, "bar", sm.GetOr("foo", ""))
		sm.SetKV("foo", "bar-changed")
		assert.Equal(t, "bar-changed", sm.GetOr("foo", ""))
		sm.Reset()
	})

	t.Run("Set", func(t *testing.T) {
		err := sm.Set(KV{})
		assert.Nil(t, err)
		err = sm.Set(KV{"foo"})
		assert.ErrorIs(t, ErrInvalidKVS, err)
		err = sm.Set(KV{"foo", "bar", "bar", "foo"})
		assert.Nil(t, err)
		assert.Equal(t, "bar", sm.GetOr("foo", ""))
		assert.Equal(t, "foo", sm.GetOr("bar", ""))
		sm.Reset()
	})

	t.Run("Get", func(t *testing.T) {
		_, ok := sm.Get("foo")
		assert.False(t, ok)
		sm.SetKV("foo", "bar")
		value, ok := sm.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "bar", value)
		sm.Reset()
	})

	t.Run("GetString", func(t *testing.T) {
		_, ok := sm.GetString("foo")
		assert.False(t, ok)
		sm.SetKV("foo", "bar")
		value, ok := sm.GetString("foo")
		assert.True(t, ok)
		assert.Equal(t, "bar", value)
		sm.Reset()
	})

	t.Run("GetStringOr", func(t *testing.T) {
		value := sm.GetStringOr("foo", "bar")
		assert.Equal(t, "bar", value)
		sm.SetKV("foo", "foo")
		value = sm.GetStringOr("foo", "bar")
		assert.Equal(t, "foo", value)
		sm.Reset()
	})

	t.Run("GetOr", func(t *testing.T) {
		value := sm.GetOr("foo", "bar")
		assert.Equal(t, "bar", value)
		sm.SetKV("foo", "foo")
		value = sm.GetOr("foo", "bar")
		assert.Equal(t, "foo", value)
		sm.Reset()
	})

	t.Run("HasKey", func(t *testing.T) {
		ok := sm.HasKey("foo")
		assert.False(t, ok)
		sm.SetKV("foo", "")
		assert.True(t, sm.HasKey("foo"))
		sm.Reset()
	})

	t.Run("HasValue", func(t *testing.T) {
		assert.False(t, sm.HasValue("bar"))
		sm.SetKV("foo", "bar")
		assert.True(t, sm.HasValue("bar"))
		sm.Reset()
	})

	t.Run("Keys", func(t *testing.T) {
		keys := sm.Keys()
		assert.Equal(t, 0, len(keys))
		expected := []string{"foo1", "foo2", "foo3"}
		for _, key := range expected {
			sm.SetKV(key, "")
		}
		keys = sm.Keys()
		var actual []string
		for _, key := range keys {
			actual = append(actual, key.(string))
		}

		assert.Equal(t, expected, actual)
		sm.Reset()
	})

	t.Run("Values", func(t *testing.T) {
		values := sm.Values()
		assert.Equal(t, 0, len(values))
		expected := []string{"foo1", "foo2", "foo3"}
		for _, key := range expected {
			sm.SetKV(key, key)
		}
		values = sm.Values()
		var actual []string
		for _, value := range values {
			actual = append(actual, value.(string))
		}

		assert.Equal(t, expected, actual)
		sm.Reset()
	})

	t.Run("Range", func(t *testing.T) {
		var keys, values []string
		sm.Range(func(key, value any) {
			keys = append(keys, key.(string))
			values = append(values, value.(string))
		})
		assert.Len(t, keys, 0)
		assert.Len(t, values, 0)

		expected := []string{"foo1", "foo2", "foo3"}
		for _, key := range expected {
			sm.SetKV(key, key)
		}
		sm.Range(func(key, value any) {
			keys = append(keys, key.(string))
			values = append(values, value.(string))
		})
		assert.Equal(t, expected, keys)
		assert.Equal(t, expected, values)
		sm.Reset()
	})

	t.Run("RangeIf", func(t *testing.T) {
		var keys, values []string
		sm.RangeIf(func(key, value any) bool {
			keys = append(keys, key.(string))
			values = append(values, value.(string))
			return true
		})
		assert.Len(t, keys, 0)
		assert.Len(t, values, 0)

		expected := []string{"foo1", "foo2", "foo3"}
		for _, key := range expected {
			sm.SetKV(key, key)
		}
		sm.RangeIf(func(key, value any) bool {
			keys = append(keys, key.(string))
			values = append(values, value.(string))
			if key.(string) == "foo1" {
				return false
			}
			return true
		})
		assert.Equal(t, []string{"foo1"}, keys)
		assert.Equal(t, []string{"foo1"}, values)
		sm.Reset()
	})

	t.Run("Remove", func(t *testing.T) {
		_, ok := sm.Remove("foo")
		assert.False(t, ok)
		sm.SetKV("foo", "bar")
		value, ok := sm.Remove("foo")
		assert.True(t, ok)
		assert.Equal(t, "bar", value)
		assert.False(t, sm.HasKey("foo"))
		assert.False(t, sm.HasValue("bar"))
		sm.Reset()
	})

	t.Run("Insert", func(t *testing.T) {
		data := New()
		data.SetKV("foo", "bar")
		sm.SetKV("foo1", "bar1")
		sm.Insert(data)
		assert.True(t, sm.HasKey("foo"))
		assert.True(t, sm.HasValue("bar"))
		sm.Reset()
	})

	t.Run("Copy", func(t *testing.T) {
		sm.SetKV("foo", "bar")
		data := sm.Copy()
		assert.True(t, data.HasKey("foo"))
		assert.True(t, data.HasValue("bar"))
		sm.SetKV("foo", "bar1")
		assert.True(t, data.HasKey("foo"))
		assert.True(t, data.HasValue("bar"))
		sm.Reset()
	})

	t.Run("Format", func(t *testing.T) {
		format := sm.Format()
		assert.Equal(t, []string{}, format)
		sm.SetKV("foo1", "bar1")
		sm.SetKV("foo2", "bar2")
		sm.SetKV("foo3", "")
		format = sm.Format()
		assert.Equal(t, []string{"foo1=bar1", "foo2=bar2", "foo3="}, format)
		sm.Reset()
	})
}
