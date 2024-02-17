package stringx

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func FuzzReplacerReplace(f *testing.F) {
	keywords := make(map[string]string)
	for i := 0; i < 20; i++ {
		keywords[Randn(rand.Intn(10)+5)] = Randn(rand.Intn(5) + 1)
	}
	rep := NewReplacer(keywords)
	printableKeywords := func() string {
		var buf strings.Builder
		for k, v := range keywords {
			fmt.Fprintf(&buf, "%q: %q,\n", k, v)
		}
		return buf.String()
	}

	f.Add(50)
	f.Fuzz(func(t *testing.T, n int) {
		text := Randn(rand.Intn(n%50+50) + 1)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("mapping: %s\ntext: %s", printableKeywords(), text)
			}
		}()
		val := rep.Replace(text)
		keys := rep.(*replacer).node.find([]rune(val))
		if len(keys) > 0 {
			t.Errorf("mapping: %s\ntext: %s\nresult: %s\nmatch: %v",
				printableKeywords(), text, val, keys)
		}
	})
}
