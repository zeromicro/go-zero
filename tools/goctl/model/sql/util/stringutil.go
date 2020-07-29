package util

import (
	"strings"
	"unicode"
)

func FormatField(field string) (snakeCase, upperCamelCase, lowerCamelCase string) {
	snakeCase = field
	list := strings.Split(field, "_")
	upperCaseList := make([]string, 0)
	lowerCaseList := make([]string, 0)
	for index, word := range list {
		upperStart := convertUpperStart(word)
		lowerStart := convertLowerStart(word)
		upperCaseList = append(upperCaseList, upperStart)
		if index == 0 {
			lowerCaseList = append(lowerCaseList, lowerStart)
		} else {
			lowerCaseList = append(lowerCaseList, upperStart)
		}
	}
	upperCamelCase = strings.Join(upperCaseList, "")
	lowerCamelCase = strings.Join(lowerCaseList, "")
	return
}

func convertLowerStart(in string) string {
	var resp []rune
	for index, r := range in {
		if index == 0 {
			r = unicode.ToLower(r)
		}
		resp = append(resp, r)
	}
	return string(resp)
}

func convertUpperStart(in string) string {
	var resp []rune
	for index, r := range in {
		if index == 0 {
			r = unicode.ToUpper(r)
		}
		resp = append(resp, r)
	}
	return string(resp)
}
