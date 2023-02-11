package stringx

import (
	"strings"
)

type (
	// Replacer interface wraps the Replace method.
	Replacer interface {
		Replace(text string) string
	}

	replacer struct {
		*node
		mapping map[string]string
	}
)

// NewReplacer returns a Replacer.
func NewReplacer(mapping map[string]string) Replacer {
	rep := &replacer{
		node:    new(node),
		mapping: mapping,
	}
	for k := range mapping {
		rep.add(k)
	}
	rep.build()

	return rep
}

// Replace replaces text with given substitutes.
func (r *replacer) Replace(text string) string {
	var buf strings.Builder
	var nextStart int
	target := []rune(text)
	cur := r.node

	for len(target) != 0 {
		used, jump, matched := cur.longestMatch(target, nextStart)
		if matched {
			replaced := r.mapping[string(target[:used])]
			target = append([]rune(replaced), target[used:]...)
			cur = r.node
			nextStart = 0
		} else {
			buf.WriteString(string(target[:used]))
			target = target[used:]
			if jump != nil {
				cur = jump
				nextStart = jump.depth
			} else {
				cur = r.node
				nextStart = 0
			}
		}
	}

	return buf.String()
}
