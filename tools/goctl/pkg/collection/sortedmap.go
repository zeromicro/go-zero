package sortedmap

import (
	"container/list"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

var (
	ErrInvalidKVExpression = errors.New(`invalid key-value expression`)
	ErrInvalidKVS          = errors.New("the length of kv must be an even number")
)

type KV []any

type SortedMap struct {
	kv   *list.List
	keys map[any]*list.Element
}

func New() *SortedMap {
	return &SortedMap{
		kv:   list.New(),
		keys: make(map[any]*list.Element),
	}
}

func (m *SortedMap) SetExpression(expression string) (key, value any, err error) {
	idx := strings.Index(expression, "=")
	if idx == -1 {
		return "", "", ErrInvalidKVExpression
	}
	key = expression[:idx]
	if len(expression) == idx {
		value = ""
	} else {
		value = expression[idx+1:]
	}
	if keys, ok := key.(string); ok && stringx.ContainsWhiteSpace(keys) {
		return "", "", ErrInvalidKVExpression
	}
	if values, ok := value.(string); ok && stringx.ContainsWhiteSpace(values) {
		return "", "", ErrInvalidKVExpression
	}
	if len(key.(string)) == 0 {
		return "", "", ErrInvalidKVExpression
	}

	m.SetKV(key, value)
	return
}

func (m *SortedMap) SetKV(key, value any) {
	e, ok := m.keys[key]
	if !ok {
		e = m.kv.PushBack(KV{
			key, value,
		})
	} else {
		e.Value.(KV)[1] = value
	}
	m.keys[key] = e
}

func (m *SortedMap) Set(kv KV) error {
	if len(kv) == 0 {
		return nil
	}
	if len(kv)%2 != 0 {
		return ErrInvalidKVS
	}
	for idx := 0; idx < len(kv); idx += 2 {
		m.SetKV(kv[idx], kv[idx+1])
	}
	return nil
}

func (m *SortedMap) Get(key any) (any, bool) {
	e, ok := m.keys[key]
	if !ok {
		return nil, false
	}
	return e.Value.(KV)[1], true
}

func (m *SortedMap) GetOr(key, dft any) any {
	e, ok := m.keys[key]
	if !ok {
		return dft
	}
	return e.Value.(KV)[1]
}

func (m *SortedMap) GetString(key any) (string, bool) {
	value, ok := m.Get(key)
	if !ok {
		return "", false
	}
	vs, ok := value.(string)
	return vs, ok
}

func (m *SortedMap) GetStringOr(key any, dft string) string {
	value, ok := m.GetString(key)
	if !ok {
		return dft
	}
	return value
}

func (m *SortedMap) HasKey(key any) bool {
	_, ok := m.keys[key]
	return ok
}

func (m *SortedMap) HasValue(value any) bool {
	var contains bool
	m.RangeIf(func(key, v any) bool {
		if value == v {
			contains = true
			return false
		}
		return true
	})
	return contains
}

func (m *SortedMap) Keys() []any {
	keys := make([]any, 0)
	next := m.kv.Front()
	for next != nil {
		keys = append(keys, next.Value.(KV)[0])
		next = next.Next()
	}
	return keys
}

func (m *SortedMap) Values() []any {
	keys := m.Keys()
	values := make([]any, len(keys))
	for idx, key := range keys {
		values[idx] = m.keys[key].Value.(KV)[1]
	}
	return values
}

func (m *SortedMap) Range(iterator func(key, value any)) {
	next := m.kv.Front()
	for next != nil {
		value := next.Value.(KV)
		iterator(value[0], value[1])
		next = next.Next()
	}
}

func (m *SortedMap) RangeIf(iterator func(key, value any) bool) {
	next := m.kv.Front()
	for next != nil {
		value := next.Value.(KV)
		loop := iterator(value[0], value[1])
		if !loop {
			return
		}
		next = next.Next()
	}
}

func (m *SortedMap) Remove(key any) (value any, ok bool) {
	v, ok := m.keys[key]
	if !ok {
		return nil, false
	}
	value = v.Value.(KV)[1]
	ok = true
	m.kv.Remove(v)
	delete(m.keys, key)
	return
}

func (m *SortedMap) Insert(sm *SortedMap) {
	sm.Range(func(key, value any) {
		m.SetKV(key, value)
	})
}

func (m *SortedMap) Copy() *SortedMap {
	sm := New()
	m.Range(func(key, value any) {
		sm.SetKV(key, value)
	})
	return sm
}

func (m *SortedMap) Format() []string {
	format := make([]string, 0)
	m.Range(func(key, value any) {
		format = append(format, fmt.Sprintf("%s=%s", key, value))
	})
	return format
}

func (m *SortedMap) Reset() {
	m.kv.Init()
	for key := range m.keys {
		delete(m.keys, key)
	}
}
