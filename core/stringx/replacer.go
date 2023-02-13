package stringx

import (
	"sort"
	"strings"
)

type (
	// Replacer interface wraps the Replace method.
	Replacer interface {
		Replace(text string) string
		Replace1(text string) string
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
func (r *replacer) Replace1(text string) string {
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

func (r *replacer) Replace(text string) string {
	for i := 0; i < 2; i++ {
		var ok bool
		if text, ok = r.doReplace(text); !ok {
			return text
		}
	}

	return text
}

func (r *replacer) doReplace(text string) (string, bool) {
	chars := []rune(text)
	scopes := r.find(chars)
	if len(scopes) == 0 {
		return text, false
	}

	sort.Slice(scopes, func(i, j int) bool {
		if scopes[i].start < scopes[j].start {
			return true
		}
		if scopes[i].start == scopes[j].start {
			return scopes[i].stop > scopes[j].stop
		}
		return false
	})

	var buf strings.Builder
	var index int
	for i := 0; i < len(scopes); i++ {
		scp := &scopes[i]
		if scp.start < index {
			continue
		}

		buf.WriteString(string(chars[index:scp.start]))
		buf.WriteString(r.mapping[string(chars[scp.start:scp.stop])])
		index = scp.stop
	}
	if index < len(chars) {
		buf.WriteString(string(chars[index:]))
	}

	return buf.String(), true
}
