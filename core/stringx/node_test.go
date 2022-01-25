package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuzzNodeCase1(t *testing.T) {
	keywords := []string{
		"cs8Zh",
		"G1OihlVuBz",
		"K6azS2FBHjI",
		"DQKvghI4",
		"l7bA86Sze",
		"tjBLZhCao",
		"nEsXmVzP",
		"cbRh8UE1nO3s",
		"Wta3R2WcbGP",
		"jpOIcA",
		"TtkRr4k9hI",
		"OKbSo0clAYTtk",
		"uJs1WToEanlKV",
		"05Y02iFD2",
		"x2uJs1WToEanlK",
		"ieaSWe",
		"Kg",
		"FD2bCKFazH",
	}
	scopes := []scope{
		{62, 72},
		{52, 65},
		{21, 34},
		{1, 10},
		{19, 33},
		{36, 42},
		{42, 44},
		{7, 17},
	}
	n := new(node)
	for _, key := range keywords {
		n.add(key)
	}
	n.build()
	assert.ElementsMatch(t, scopes, n.find([]rune("Z05Y02iFD2bCKFazHtrx2uJs1WToEanlKVWKieaSWeKgmnUXV0ZjOKbSo0clAYTtkRr4k9hI")))
}

func TestFuzzNodeCase2(t *testing.T) {
	keywords := []string{
		"IP1NPsJKIvt",
		"Iw7hQARwSTw",
		"SmZIcA",
		"OyxHPYkoQzFO",
		"3suCnuSAS5d",
		"HUMpbi",
		"HPdvbGGpY",
		"49qjMtR8bEt",
		"a0zrrGKVTJ2",
		"WbOBcszeo1FfZ",
		"8tHUi5PJI",
		"Oa2Di",
		"6ZWa5nr1tU",
		"o0LJRfmeXB9bF9",
		"veF0ehKxH",
		"Qp73r",
		"B6Rmds4ELY8",
		"uNpGybQZG",
		"Ogm3JqicRZlA4n",
		"FL6LVErKomc84H",
		"qv2Pi0xJj3cR1",
		"bPWLBg4",
		"hYN8Q4M1sw",
		"ExkTgNklmlIx",
		"eVgHHDOxOUEj",
		"5WPEVv0tR",
		"CPjnOAqUZgV",
		"oR3Ogtz",
		"jwk1Zbg",
		"DYqguyk8h",
		"rieroDmpvYFK",
		"MQ9hZnMjDqrNQe",
		"EhM4KqkCBd",
		"m9xalj6q",
		"d5CTL5mzK",
		"XJOoTvFtI8U",
		"iFAwspJ",
		"iGv8ErnRZIuSWX",
		"i8C1BqsYX",
		"vXN1KOaOgU",
		"GHJFB",
		"Y6OlAqbZxYG8",
		"dzd4QscSih4u",
		"SsLYMkKvB9APx",
		"gi0huB3",
		"CMICHDCSvSrgiACXVkN",
		"MwOvyHbaxdaqpZpU",
		"wOvyHbaxdaqpZpUbI",
		"2TT5WEy",
		"eoCq0T2MC",
		"ZpUbI7",
		"oCq0T2MCp",
		"CpLFgLg0g",
		"FgLg0gh",
		"w5awC5HeoCq",
		"1c",
	}
	scopes := []scope{
		{0, 19},
		{57, 73},
		{58, 75},
		{47, 54},
		{29, 38},
		{70, 76},
		{30, 39},
		{37, 46},
		{40, 47},
		{22, 33},
		{92, 94},
	}
	n := new(node)
	for _, key := range keywords {
		n.add(key)
	}
	n.build()
	assert.ElementsMatch(t, scopes, n.find([]rune("CMICHDCSvSrgiACXVkNF9lw5awC5HeoCq0T2MCpLFgLg0gh2TT5WEyINrMwOvyHbaxdaqpZpUbI7SpIY5yVWf33MuX7K1c")))
}

func TestFuzzNodeCase3(t *testing.T) {
	keywords := []string{
		"QAraACKOftI4",
		"unRmd2EO0",
		"s25OtuoU",
		"aGlmn7KnbE4HCX",
		"kuK6Uh",
		"ckuK6Uh",
		"uK6Uh",
		"Iy",
		"h",
		"PMSSUNvyi",
		"ahz0i",
		"Lhs4XZ1e",
		"shPp1Va7aQNVme",
		"yIUckuK6Uh",
		"pKjIyI",
		"jIyIUckuK6Uh",
		"UckuK6Uh",
		"Uh",
		"JPAULjQgHJ",
		"Wp",
		"sbkZxXurrI",
		"pKjIyIUckuK6Uh",
	}
	scopes := []scope{
		{9, 15},
		{8, 15},
		{5, 15},
		{1, 7},
		{10, 15},
		{3, 15},
		{0, 2},
		{1, 15},
		{7, 15},
		{13, 15},
		{4, 6},
		{14, 15},
	}
	n := new(node)
	for _, key := range keywords {
		n.add(key)
	}
	n.build()
	assert.ElementsMatch(t, scopes, n.find([]rune("WpKjIyIUckuK6Uh")))
}

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
