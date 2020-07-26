package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAddrs(t *testing.T) {
	tests := []struct {
		addrs  []string
		expect string
	}{
		{
			addrs:  []string{"a", "b"},
			expect: "a,b",
		},
		{
			addrs:  []string{"a", "b", "c"},
			expect: "a,b,c",
		},
		{
			addrs:  []string{},
			expect: "",
		},
		{
			addrs:  nil,
			expect: "",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expect, FormatAddr(test.addrs))
	}
}
