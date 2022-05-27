package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacer_Replace(t *testing.T) {
	mapping := map[string]string{
		"一二三四": "1234",
		"二三":   "23",
		"二":    "2",
	}
	assert.Equal(t, "零1234五", NewReplacer(mapping).Replace("零一二三四五"))
}

func TestReplacer_ReplaceOverlap(t *testing.T) {
	mapping := map[string]string{
		"3d": "34",
		"bc": "23",
	}
	assert.Equal(t, "a234e", NewReplacer(mapping).Replace("abcde"))
}

func TestReplacer_ReplaceSingleChar(t *testing.T) {
	mapping := map[string]string{
		"二": "2",
	}
	assert.Equal(t, "零一2三四五", NewReplacer(mapping).Replace("零一二三四五"))
}

func TestReplacer_ReplaceExceedRange(t *testing.T) {
	mapping := map[string]string{
		"二三四五六": "23456",
	}
	assert.Equal(t, "零一二三四五", NewReplacer(mapping).Replace("零一二三四五"))
}

func TestReplacer_ReplacePartialMatch(t *testing.T) {
	mapping := map[string]string{
		"二三四七": "2347",
	}
	assert.Equal(t, "零一二三四五", NewReplacer(mapping).Replace("零一二三四五"))
}

func TestReplacer_ReplaceMultiMatches(t *testing.T) {
	mapping := map[string]string{
		"二三": "23",
	}
	assert.Equal(t, "零一23四五一23四五", NewReplacer(mapping).Replace("零一二三四五一二三四五"))
}

func TestReplacer_ReplaceJumpToFail(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"cde":  "234",
	}
	assert.Equal(t, "ab234fg", NewReplacer(mapping).Replace("abcdefg"))
}

func TestReplacer_ReplaceJumpToFailDup(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"ccde": "2234",
	}
	assert.Equal(t, "ab2234fg", NewReplacer(mapping).Replace("abccdefg"))
}

func TestReplacer_ReplaceJumpToFailEnding(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"cdef": "2345",
	}
	assert.Equal(t, "ab2345", NewReplacer(mapping).Replace("abcdef"))
}

func TestReplacer_ReplaceEmpty(t *testing.T) {
	mapping := map[string]string{
		"bcdf": "1235",
		"cdef": "2345",
	}
	assert.Equal(t, "", NewReplacer(mapping).Replace(""))
}

func TestFuzzReplacerCase1(t *testing.T) {
	keywords := map[string]string{
		"yQyJykiqoh":     "xw",
		"tgN70z":         "Q2P",
		"tXKhEn":         "w1G8",
		"5nfOW1XZO":      "GN",
		"f4Ov9i9nHD":     "cT",
		"1ov9Q":          "Y",
		"7IrC9n":         "400i",
		"JQLxonpHkOjv":   "XI",
		"DyHQ3c7":        "Ygxux",
		"ffyqJi":         "u",
		"UHuvXrbD8pni":   "dN",
		"LIDzNbUlTX":     "g",
		"yN9WZh2rkc8Q":   "3U",
		"Vhk11rz8CObceC": "jf",
		"R0Rt4H2qChUQf":  "7U5M",
		"MGQzzPCVKjV9":   "yYz",
		"B5jUUl0u1XOY":   "l4PZ",
		"pdvp2qfLgG8X":   "BM562",
		"ZKl9qdApXJ2":    "T",
		"37jnugkSevU66":  "aOHFX",
	}
	rep := NewReplacer(keywords)
	text := "yjF8fyqJiiqrczOCVyoYbLvrMpnkj"
	val := rep.Replace(text)
	keys := rep.(*replacer).node.find([]rune(val))
	if len(keys) > 0 {
		t.Errorf("result: %s, match: %v", val, keys)
	}
}

func TestFuzzReplacerCase2(t *testing.T) {
	keywords := map[string]string{
		"dmv2SGZvq9Yz":   "TE",
		"rCL5DRI9uFP8":   "hvsc8",
		"7pSA2jaomgg":    "v",
		"kWSQvjVOIAxR":   "Oje",
		"hgU5bYYkD3r6":   "qCXu",
		"0eh6uI":         "MMlt",
		"3USZSl85EKeMzw": "Pc",
		"JONmQSuXa":      "dX",
		"EO1WIF":         "G",
		"uUmFJGVmacjF":   "1N",
		"DHpw7":          "M",
		"NYB2bm":         "CPya",
		"9FiNvBAHHNku5":  "7FlDE",
		"tJi3I4WxcY":     "q5",
		"sNJ8Z1ToBV0O":   "tl",
		"0iOg72QcPo":     "RP",
		"pSEqeL":         "5KZ",
		"GOyYqTgmvQ":     "9",
		"Qv4qCsj":        "nl52E",
		"wNQ5tOutYu5s8":  "6iGa",
	}
	rep := NewReplacer(keywords)
	text := "AoRxrdKWsGhFpXwVqMLWRL74OukwjBuBh0g7pSrk"
	val := rep.Replace(text)
	keys := rep.(*replacer).node.find([]rune(val))
	if len(keys) > 0 {
		t.Errorf("result: %s, match: %v", val, keys)
	}
}
