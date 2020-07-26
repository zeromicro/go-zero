package filter

import (
	"strings"

	"zero/stash/config"

	"github.com/globalsign/mgo/bson"
)

const (
	filterDrop         = "drop"
	filterRemoveFields = "remove_field"
	opAnd              = "and"
	opOr               = "or"
	typeContains       = "contains"
	typeMatch          = "match"
)

type FilterFunc func(map[string]interface{}) map[string]interface{}

func CreateFilters(c config.Config) []FilterFunc {
	var filters []FilterFunc

	for _, f := range c.Filters {
		switch f.Action {
		case filterDrop:
			filters = append(filters, DropFilter(f.Conditions))
		case filterRemoveFields:
			filters = append(filters, RemoveFieldFilter(f.Fields))
		}
	}

	return filters
}

func DropFilter(conds []config.Condition) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		var qualify bool
		for _, cond := range conds {
			var qualifyOnce bool
			switch cond.Type {
			case typeMatch:
				qualifyOnce = cond.Value == m[cond.Key]
			case typeContains:
				if val, ok := m[cond.Key].(string); ok {
					qualifyOnce = strings.Contains(val, cond.Value)
				}
			}

			switch cond.Op {
			case opAnd:
				if !qualifyOnce {
					return m
				} else {
					qualify = true
				}
			case opOr:
				if qualifyOnce {
					qualify = true
				}
			}
		}

		if qualify {
			return nil
		} else {
			return m
		}
	}
}

func RemoveFieldFilter(fields []string) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		for _, field := range fields {
			delete(m, field)
		}
		return m
	}
}

func AddUriFieldFilter(inField, outFirld string) FilterFunc {
	return func(m map[string]interface{}) map[string]interface{} {
		if val, ok := m[inField].(string); ok {
			var datas []string
			idx := strings.Index(val, "?")
			if idx < 0 {
				datas = strings.Split(val, "/")
			} else {
				datas = strings.Split(val[:idx], "/")
			}

			for i, data := range datas {
				if bson.IsObjectIdHex(data) {
					datas[i] = "*"
				}
			}

			m[outFirld] = strings.Join(datas, "/")
		}

		return m
	}
}
