package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie(t *testing.T) {
	tests := []struct {
		input    string
		output   string
		keywords []string
		found    bool
	}{
		{
			input:  "日本AV演员兼电视、电影演员。苍井空AV女优是xx出道, 日本AV女优们最精彩的表演是AV演员色情表演",
			output: "日本****兼电视、电影演员。*****女优是xx出道, ******们最精彩的表演是******表演",
			keywords: []string{
				"AV演员",
				"苍井空",
				"AV",
				"日本AV女优",
				"AV演员色情",
			},
			found: true,
		},
		{
			input:    "完全和谐的文本完全和谐的文本",
			output:   "完全和谐的文本完全和谐的文本",
			keywords: nil,
			found:    false,
		},
		{
			input:  "就一个字不对",
			output: "就*个字不对",
			keywords: []string{
				"一",
			},
			found: true,
		},
		{
			input:  "就一对, AV",
			output: "就*对, **",
			keywords: []string{
				"一",
				"AV",
			},
			found: true,
		},
		{
			input:  "就一不对, AV",
			output: "就**对, **",
			keywords: []string{
				"一",
				"一不",
				"AV",
			},
			found: true,
		},
		{
			input:  "就对, AV",
			output: "就对, **",
			keywords: []string{
				"AV",
			},
			found: true,
		},
		{
			input:  "就对, 一不",
			output: "就对, **",
			keywords: []string{
				"一",
				"一不",
			},
			found: true,
		},
		{
			input:    "",
			output:   "",
			keywords: nil,
			found:    false,
		},
	}

	trie := NewTrie([]string{
		"", // no hurts for empty keywords
		"一",
		"一不",
		"AV",
		"AV演员",
		"苍井空",
		"AV演员色情",
		"日本AV女优",
	})

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output, keywords, ok := trie.Filter(test.input)
			assert.Equal(t, test.found, ok)
			assert.Equal(t, test.output, output)
			assert.ElementsMatch(t, test.keywords, keywords)
			keywords = trie.FindKeywords(test.input)
			assert.ElementsMatch(t, test.keywords, keywords)
		})
	}
}

func TestTrieSingleWord(t *testing.T) {
	trie := NewTrie([]string{
		"闹",
	}, WithMask('#'))
	output, keywords, ok := trie.Filter("今晚真热闹")
	assert.ElementsMatch(t, []string{"闹"}, keywords)
	assert.True(t, ok)
	assert.Equal(t, "今晚真热#", output)
}

func TestTrieOverlap(t *testing.T) {
	trie := NewTrie([]string{
		"一二三四五",
		"二三四五六七八",
	}, WithMask('#'))
	output, keywords, ok := trie.Filter("零一二三四五六七八九十")
	assert.ElementsMatch(t, []string{
		"一二三四五",
		"二三四五六七八",
	}, keywords)
	assert.True(t, ok)
	assert.Equal(t, "零########九十", output)
}

func TestTrieNested(t *testing.T) {
	trie := NewTrie([]string{
		"一二三",
		"一二三四五",
		"一二三四五六七八",
	}, WithMask('#'))
	output, keywords, ok := trie.Filter("零一二三四五六七八九十")
	assert.ElementsMatch(t, []string{
		"一二三",
		"一二三四五",
		"一二三四五六七八",
	}, keywords)
	assert.True(t, ok)
	assert.Equal(t, "零########九十", output)
}

func BenchmarkTrie(b *testing.B) {
	b.ReportAllocs()

	trie := NewTrie([]string{
		"A",
		"AV",
		"AV演员",
		"苍井空",
		"AV演员色情",
		"日本AV女优",
	})

	for i := 0; i < b.N; i++ {
		trie.Filter("日本AV演员兼电视、电影演员。苍井空AV女优是xx出道, 日本AV女优们最精彩的表演是AV演员色情表演")
	}
}
