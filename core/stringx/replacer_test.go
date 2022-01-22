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
