package collection

import (
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/logx"
)

const (
	unmanaged = iota
	untyped
	intType
	int64Type
	uintType
	uint64Type
	stringType
)

type Set struct {
	data map[interface{}]lang.PlaceholderType
	tp   int
}

func NewSet() *Set {
	return &Set{
		data: make(map[interface{}]lang.PlaceholderType),
		tp:   untyped,
	}
}

func NewUnmanagedSet() *Set {
	return &Set{
		data: make(map[interface{}]lang.PlaceholderType),
		tp:   unmanaged,
	}
}

func (s *Set) Add(i ...interface{}) {
	for _, each := range i {
		s.add(each)
	}
}

func (s *Set) AddInt(ii ...int) {
	for _, each := range ii {
		s.add(each)
	}
}

func (s *Set) AddInt64(ii ...int64) {
	for _, each := range ii {
		s.add(each)
	}
}

func (s *Set) AddUint(ii ...uint) {
	for _, each := range ii {
		s.add(each)
	}
}

func (s *Set) AddUint64(ii ...uint64) {
	for _, each := range ii {
		s.add(each)
	}
}

func (s *Set) AddStr(ss ...string) {
	for _, each := range ss {
		s.add(each)
	}
}

func (s *Set) Contains(i interface{}) bool {
	if len(s.data) == 0 {
		return false
	}

	s.validate(i)
	_, ok := s.data[i]
	return ok
}

func (s *Set) Keys() []interface{} {
	var keys []interface{}

	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

func (s *Set) KeysInt() []int {
	var keys []int

	for key := range s.data {
		if intKey, ok := key.(int); !ok {
			continue
		} else {
			keys = append(keys, intKey)
		}
	}

	return keys
}

func (s *Set) KeysInt64() []int64 {
	var keys []int64

	for key := range s.data {
		if intKey, ok := key.(int64); !ok {
			continue
		} else {
			keys = append(keys, intKey)
		}
	}

	return keys
}

func (s *Set) KeysUint() []uint {
	var keys []uint

	for key := range s.data {
		if intKey, ok := key.(uint); !ok {
			continue
		} else {
			keys = append(keys, intKey)
		}
	}

	return keys
}

func (s *Set) KeysUint64() []uint64 {
	var keys []uint64

	for key := range s.data {
		if intKey, ok := key.(uint64); !ok {
			continue
		} else {
			keys = append(keys, intKey)
		}
	}

	return keys
}

func (s *Set) KeysStr() []string {
	var keys []string

	for key := range s.data {
		if strKey, ok := key.(string); !ok {
			continue
		} else {
			keys = append(keys, strKey)
		}
	}

	return keys
}

func (s *Set) Remove(i interface{}) {
	s.validate(i)
	delete(s.data, i)
}

func (s *Set) Count() int {
	return len(s.data)
}

func (s *Set) add(i interface{}) {
	switch s.tp {
	case unmanaged:
		// do nothing
	case untyped:
		s.setType(i)
	default:
		s.validate(i)
	}
	s.data[i] = lang.Placeholder
}

func (s *Set) setType(i interface{}) {
	if s.tp != untyped {
		return
	}

	switch i.(type) {
	case int:
		s.tp = intType
	case int64:
		s.tp = int64Type
	case uint:
		s.tp = uintType
	case uint64:
		s.tp = uint64Type
	case string:
		s.tp = stringType
	}
}

func (s *Set) validate(i interface{}) {
	if s.tp == unmanaged {
		return
	}

	switch i.(type) {
	case int:
		if s.tp != intType {
			logx.Errorf("Error: element is int, but set contains elements with type %d", s.tp)
		}
	case int64:
		if s.tp != int64Type {
			logx.Errorf("Error: element is int64, but set contains elements with type %d", s.tp)
		}
	case uint:
		if s.tp != uintType {
			logx.Errorf("Error: element is uint, but set contains elements with type %d", s.tp)
		}
	case uint64:
		if s.tp != uint64Type {
			logx.Errorf("Error: element is uint64, but set contains elements with type %d", s.tp)
		}
	case string:
		if s.tp != stringType {
			logx.Errorf("Error: element is string, but set contains elements with type %d", s.tp)
		}
	}
}
