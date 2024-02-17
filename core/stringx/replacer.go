package stringx

import (
	"sort"
	"strings"
)

// replace more than once to avoid overlapped keywords after replace.
// only try 2 times to avoid too many or infinite loops.
const replaceTimes = 2

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
	for i := 0; i < replaceTimes; i++ {
		var replaced bool
		if text, replaced = r.doReplace(text); !replaced {
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
