package stringx

import "strings"

type (
	// Replacer interface wraps the Replace method.
	Replacer interface {
		Replace(text string) string
	}

	replacer struct {
		node
		mapping map[string]string
	}
)

// NewReplacer returns a Replacer.
func NewReplacer(mapping map[string]string) Replacer {
	rep := &replacer{
		mapping: mapping,
	}
	for k := range mapping {
		rep.add(k)
	}

	return rep
}

func (r *replacer) Replace(text string) string {
	var builder strings.Builder
	chars := []rune(text)
	size := len(chars)
	start := -1

	for i := 0; i < size; i++ {
		child, ok := r.children[chars[i]]
		if !ok {
			builder.WriteRune(chars[i])
			continue
		}

		if start < 0 {
			start = i
		}
		end := -1
		if child.end {
			end = i + 1
		}

		j := i + 1
		for ; j < size; j++ {
			grandchild, ok := child.children[chars[j]]
			if !ok {
				break
			}

			child = grandchild
			if child.end {
				end = j + 1
				i = j
			}
		}

		if end > 0 {
			i = j - 1
			builder.WriteString(r.mapping[string(chars[start:end])])
		} else {
			builder.WriteRune(chars[i])
		}
		start = -1
	}

	return builder.String()
}
