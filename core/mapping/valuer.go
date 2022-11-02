package mapping

type (
	// A Valuer interface defines the way to get values from the underlying object with keys.
	Valuer interface {
		// Value gets the value associated with the given key.
		Value(key string) (interface{}, bool)
	}

	// A ValuerWithParent defines a node that has a parent node.
	ValuerWithParent interface {
		Valuer
		Parent() ValuerWithParent
	}

	// A node is a map that can use Value method to get values with given keys.
	node struct {
		current Valuer
		parent  ValuerWithParent
	}

	mapValuer       map[string]interface{}
	simpleValuer    node
	recursiveValuer node
)

func (mv mapValuer) Value(key string) (interface{}, bool) {
	v, ok := mv[key]
	return v, ok
}

// Value gets the value associated with the given key from mv.
func (sv simpleValuer) Value(key string) (interface{}, bool) {
	v, ok := sv.current.Value(key)
	return v, ok
}

func (sv simpleValuer) Parent() ValuerWithParent {
	if sv.parent == nil {
		return nil
	}

	return recursiveValuer{
		current: sv.parent,
		parent:  sv.parent.Parent(),
	}
}

// Value gets the value associated with the given key from mv,
// and it will inherit the value from parent nodes.
func (rv recursiveValuer) Value(key string) (interface{}, bool) {
	if v, ok := rv.current.Value(key); ok {
		return v, ok
	}

	if parent := rv.Parent(); parent != nil {
		return parent.Value(key)
	}

	return nil, false
}

func (rv recursiveValuer) Parent() ValuerWithParent {
	if rv.parent == nil {
		return nil
	}

	return recursiveValuer{
		current: rv.parent,
		parent:  rv.parent.Parent(),
	}
}
