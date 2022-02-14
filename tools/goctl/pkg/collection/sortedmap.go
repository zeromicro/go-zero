package kv

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidKVExpression = errors.New(`invalid key-value expression`)
var ErrInvalidKVS = errors.New("the length of kv must be a even number")

type KV []interface{}

type SortedMap struct {
	kv   *list.List
	keys map[interface{}]*list.Element
}

func NewSortedMap() *SortedMap {
	return &SortedMap{
		kv:   list.New(),
		keys: make(map[interface{}]*list.Element),
	}
}

func (m *SortedMap) SetExpression(expression string) (key interface{}, value interface{}, err error) {
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

	m.SetKV(key, value)
	return
}

func (m *SortedMap) SetKV(key, value interface{}) {
	e, ok := m.keys[key]
	if !ok {
		e = &list.Element{Value: KV{
			key, value,
		}}
		m.kv.PushBack(e)
	} else {
		e.Value.(KV)[1] = value
	}
	m.keys[key] = e
}

func (m *SortedMap) Set(kvs ...KV) error {
	if len(kvs)%2 != 0 {
		return ErrInvalidKVS
	}
	for _, kv := range kvs {
		m.SetKV(kv[0], kv[1])
	}
	return nil
}

func (m *SortedMap) Get(key interface{}) (interface{}, bool) {
	e, ok := m.keys[key]
	if !ok {
		return nil, false
	}
	return e.Value.(KV)[1], true
}

func (m *SortedMap) GetOr(key interface{}, dft interface{}) interface{} {
	e, ok := m.keys[key]
	if !ok {
		return dft
	}
	return e.Value.(KV)[1]
}

func (m *SortedMap) HasKey(key interface{}) bool {
	_, ok := m.keys[key]
	return ok
}

func (m *SortedMap) HasValue(value interface{}) bool {
	var contains bool
	m.RangeIf(func(key, v interface{}) bool {
		if value == v {
			contains = true
			return false
		}
		return true
	})
	return contains
}

func (m *SortedMap) Keys() []interface{} {
	keys := make([]interface{}, len(m.keys))
	next := m.kv.Front()
	for next != nil {
		keys = append(keys, next.Value.(KV)[0])
		next = next.Next()
	}
	return keys
}

func (m *SortedMap) Values() []interface{} {
	keys := m.Keys()
	values := make([]interface{}, len(keys))
	for _, key := range keys {
		values = append(values, m.keys[key].Value.(KV)[1])
	}
	return values
}

func (m *SortedMap) Range(iterator func(key, value interface{})) {
	next := m.kv.Front()
	for next != nil {
		value := next.Value.(KV)
		iterator(value[0], value[1])
		next = next.Next()
	}
}

func (m *SortedMap) RangeIf(iterator func(key, value interface{}) bool) {
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

func (m *SortedMap) Remove(key interface{}) (value interface{}, ok bool) {
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
	sm.Range(func(key, value interface{}) {
		m.SetKV(key, value)
	})
}

func (m *SortedMap) Copy() *SortedMap {
	sm := NewSortedMap()
	m.Range(func(key, value interface{}) {
		sm.SetKV(key, value)
	})
	return sm
}

func (m *SortedMap) Format() []string {
	var format = make([]string, len(m.keys))
	m.Range(func(key, value interface{}) {
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
