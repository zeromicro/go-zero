package mapping

type (
	Valuer interface {
		Value(key string) (interface{}, bool)
	}

	MapValuer map[string]interface{}
)

func (mv MapValuer) Value(key string) (interface{}, bool) {
	v, ok := mv[key]
	return v, ok
}
