package stringx

import "github.com/tal-tech/go-zero/core/lang"

const defaultMask = '*'

type (
	// TrieOption defines the method to customize a Trie.
	TrieOption func(trie *trieNode)

	// A Trie is a tree implementation that used to find elements rapidly.
	Trie interface {
		Filter(text string) (string, []string, bool)
		FindKeywords(text string) []string
	}

	trieNode struct {
		node
		mask rune
	}

	scope struct {
		start int
		stop  int
	}
)

// NewTrie returns a Trie.
func NewTrie(words []string, opts ...TrieOption) Trie {
	n := new(trieNode)

	for _, opt := range opts {
		opt(n)
	}
	if n.mask == 0 {
		n.mask = defaultMask
	}
	for _, word := range words {
		n.add(word)
	}

	return n
}

func (n *trieNode) Filter(text string) (sentence string, keywords []string, found bool) {
	chars := []rune(text)
	if len(chars) == 0 {
		return text, nil, false
	}

	scopes := n.findKeywordScopes(chars)
	keywords = n.collectKeywords(chars, scopes)

	for _, match := range scopes {
		// we don't care about overlaps, not bringing a performance improvement
		n.replaceWithAsterisk(chars, match.start, match.stop)
	}

	return string(chars), keywords, len(keywords) > 0
}

func (n *trieNode) FindKeywords(text string) []string {
	chars := []rune(text)
	if len(chars) == 0 {
		return nil
	}

	scopes := n.findKeywordScopes(chars)
	return n.collectKeywords(chars, scopes)
}

func (n *trieNode) collectKeywords(chars []rune, scopes []scope) []string {
	set := make(map[string]lang.PlaceholderType)
	for _, v := range scopes {
		set[string(chars[v.start:v.stop])] = lang.Placeholder
	}

	var i int
	keywords := make([]string, len(set))
	for k := range set {
		keywords[i] = k
		i++
	}

	return keywords
}

func (n *trieNode) findKeywordScopes(chars []rune) []scope {
	var scopes []scope
	size := len(chars)
	start := -1

	for i := 0; i < size; i++ {
		child, ok := n.children[chars[i]]
		if !ok {
			continue
		}

		if start < 0 {
			start = i
		}
		if child.end {
			scopes = append(scopes, scope{
				start: start,
				stop:  i + 1,
			})
		}

		for j := i + 1; j < size; j++ {
			grandchild, ok := child.children[chars[j]]
			if !ok {
				break
			}

			child = grandchild
			if child.end {
				scopes = append(scopes, scope{
					start: start,
					stop:  j + 1,
				})
			}
		}

		start = -1
	}

	return scopes
}

func (n *trieNode) replaceWithAsterisk(chars []rune, start, stop int) {
	for i := start; i < stop; i++ {
		chars[i] = n.mask
	}
}

// WithMask customizes a Trie with keywords masked as given mask char.
func WithMask(mask rune) TrieOption {
	return func(n *trieNode) {
		n.mask = mask
	}
}
