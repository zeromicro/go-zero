package stringx

import "testing"

func BenchmarkNodeFind(b *testing.B) {
	b.ReportAllocs()

	keywords := []string{
		"A",
		"AV",
		"AV演员",
		"无名氏",
		"AV演员色情",
		"日本AV女优",
	}
	trie := new(node)
	for _, keyword := range keywords {
		trie.add(keyword)
	}
	trie.build()

	for i := 0; i < b.N; i++ {
		trie.find([]rune("日本AV演员兼电视、电影演员。无名氏AV女优是xx出道, 日本AV女优们最精彩的表演是AV演员色情表演"))
	}
}
