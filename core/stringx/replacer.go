package stringx

import "strings"

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
	target := []rune(text)
	cur := r.node
	nextStart := 0
	for len(target) != 0 {
		uselessLen, matchLen, jump := cur.longestMatch(target, nextStart)
		if uselessLen > 0 {
			buf.WriteString(string(target[:uselessLen]))
			target = target[uselessLen:]
		}
		if matchLen > 0 {
			replaced := r.mapping[string(target[:matchLen])]
			target = append([]rune(replaced), target[matchLen:]...)
		}
		if jump != nil {
			cur = jump
			nextStart = jump.depth
		} else {
			cur = r.node
			nextStart = 0
		}
	}
	return buf.String()
}
