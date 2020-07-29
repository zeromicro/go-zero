package util

import (
	"bytes"
	"strings"
)

// 简单的下划线转驼峰格式
func FmtUnderLine2Camel(in string, upperStart bool) string {
	if strings.TrimSpace(in) == "" {
		return ""
	}
	var words []string
	if strings.Contains(in, "_") {
		words = strings.Split(in, "_")
		if len(words) == 0 {
			return ""
		}
	}
	if len(words) == 0 {
		return strings.Title(in)
	}
	var buffer bytes.Buffer
	for index, word := range words {
		if strings.TrimSpace(word) == "" {
			continue
		}
		bts := []byte(word)
		if index == 0 && !upperStart {
			bts[0] = bytes.ToLower([]byte{bts[0]})[0]
			buffer.Write(bts)
			continue
		}
		bts = bytes.Title(bts)
		buffer.Write(bts)
	}
	return buffer.String()
}

func Abbr(in string) string {
	ar := strings.Split(in, "_")
	var abbrBuffer bytes.Buffer
	for _, item := range ar {
		bts := []byte(item)
		if len(bts) == 0 {
			continue
		}
		abbrBuffer.Write([]byte{bts[0]})
	}
	return abbrBuffer.String()
}
