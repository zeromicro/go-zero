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
	var paths []*node
	target := []rune(text)
	cur := r.node

	for len(target) != 0 {
		uselessLen, matchLen, nextPaths := cur.longestMatch(target, paths)
		if uselessLen > 0 {
			buf.WriteString(string(target[:uselessLen]))
			target = target[uselessLen:]
		}
		if matchLen > 0 {
			replaced := r.mapping[string(target[:matchLen])]
			target = append([]rune(replaced), target[matchLen:]...)
		}
		if len(nextPaths) != 0 {
			cur = nextPaths[len(nextPaths)-1]
			paths = nextPaths
		} else {
			cur = r.node
			paths = nil
		}
	}

	return buf.String()
}
