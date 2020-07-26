package collection

type Ring struct {
	elements []interface{}
	index    int
}

func NewRing(n int) *Ring {
	return &Ring{
		elements: make([]interface{}, n),
	}
}

func (r *Ring) Add(v interface{}) {
	r.elements[r.index%len(r.elements)] = v
	r.index++
}

func (r *Ring) Take() []interface{} {
	var size int
	var start int
	if r.index > len(r.elements) {
		size = len(r.elements)
		start = r.index % len(r.elements)
	} else {
		size = r.index
	}

	elements := make([]interface{}, size)
	for i := 0; i < size; i++ {
		elements[i] = r.elements[(start+i)%len(r.elements)]
	}

	return elements
}
