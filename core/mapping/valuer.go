package mapping

type (
	// A Valuer interface defines the way to get values from the underlying object with keys.
	Valuer interface {
		// Value gets the value associated with the given key.
		Value(key string) (interface{}, bool)
	}

	// A MapValuer is a map that can use Value method to get values with given keys.
	MapValuer map[string]interface{}
)

// Value gets the value associated with the given key from mv.
func (mv MapValuer) Value(key string) (interface{}, bool) {
	v, ok := mv[key]
	return v, ok
}
