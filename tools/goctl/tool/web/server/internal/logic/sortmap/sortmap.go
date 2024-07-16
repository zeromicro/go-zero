package sortmap

import "sort"

type (
	KVPair struct {
		Key   string
		Value any
	}
	SortMap struct {
		list []*KVPair
		m    map[string]any
	}
)

func From(v map[string]any) *SortMap {
	return &SortMap{
		list: toKVPair(v),
		m:    v,
	}
}

func (m *SortMap) Range(fn func(idx int, key string, value any) error) error {
	for idx, kv := range m.list {
		if fn != nil {
			err := fn(idx, kv.Key, kv.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *SortMap) Del(key string) {
	for i, kv := range m.list {
		if kv.Key == key {
			m.list = append(m.list[:i], m.list[i+1:]...)
			return
		}
	}
}

func (m *SortMap) Get(key string) (any, bool) {
	val, ok := m.m[key]
	return val, ok
}

func toKVPair(v map[string]any) []*KVPair {
	var result []*KVPair
	for k, v := range v {
		result = append(result, &KVPair{
			Key:   k,
			Value: v,
		})
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})
	return result
}
