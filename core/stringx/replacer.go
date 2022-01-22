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
	var builder strings.Builder
	var start int
	chars := []rune(text)
	size := len(chars)

	for start < size {
		cur := r.node

		if start > 0 {
			builder.WriteString(string(chars[:start]))
		}

		for i := start; i < size; i++ {
			child, ok := cur.children[chars[i]]
			if ok {
				cur = child
			} else if cur == r.node {
				builder.WriteRune(chars[i])
				start = i + 1
				continue
			} else {
				curDepth := cur.depth
				cur = cur.fail
				child, ok = cur.children[chars[i]]
				if !ok {
					builder.WriteString(string(chars[start : i+1]))
					start = i + 1
					continue
				}

				failDepth := cur.depth
				builder.WriteString(string(chars[start : start+curDepth-failDepth]))
				cur = child
			}

			if cur.end {
				val := string(chars[i+1-cur.depth : i+1])
				builder.WriteString(r.mapping[val])
				builder.WriteString(string(chars[i+1:]))
				if i+1 >= size {
					return builder.String()
				}

				chars = []rune(builder.String())
				size = len(chars)
				builder.Reset()
				break
			}
		}

		if !cur.end {
			builder.WriteString(string(chars[start:]))
			return builder.String()
		}
	}

	return string(chars)
}
