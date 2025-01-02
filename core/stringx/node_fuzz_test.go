package stringx

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func FuzzNodeFind(f *testing.F) {
	rand.NewSource(time.Now().UnixNano())

	f.Add(10)
	f.Fuzz(func(t *testing.T, keys int) {
		str := Randn(rand.Intn(100) + 50)
		keywords := make(map[string]struct{})
		for i := 0; i < keys; i++ {
			keyword := Randn(rand.Intn(10) + 5)
			if !strings.Contains(str, keyword) {
				keywords[keyword] = struct{}{}
			}
		}

		size := len(str)
		var scopes []scope
		var n node
		for i := 0; i < size%20; i++ {
			start := rand.Intn(size)
			stop := start + rand.Intn(20) + 1
			if stop > size {
				stop = size
			}
			if start == stop {
				continue
			}

			keyword := str[start:stop]
			if _, ok := keywords[keyword]; ok {
				continue
			}

			keywords[keyword] = struct{}{}
			var pos int
			for pos <= len(str)-len(keyword) {
				val := str[pos:]
				p := strings.Index(val, keyword)
				if p < 0 {
					break
				}

				scopes = append(scopes, scope{
					start: pos + p,
					stop:  pos + p + len(keyword),
				})
				pos += p + 1
			}
		}

		for keyword := range keywords {
			n.add(keyword)
		}
		n.build()

		var buf strings.Builder
		buf.WriteString("keywords:\n")
		for key := range keywords {
			fmt.Fprintf(&buf, "\t%q,\n", key)
		}
		buf.WriteString("scopes:\n")
		for _, scp := range scopes {
			fmt.Fprintf(&buf, "\t{%d, %d},\n", scp.start, scp.stop)
		}
		fmt.Fprintf(&buf, "text:\n\t%s\n", str)
		defer func() {
			if r := recover(); r != nil {
				t.Error(buf.String())
			}
		}()
		assert.ElementsMatchf(t, scopes, n.find([]rune(str)), buf.String())
	})
}
