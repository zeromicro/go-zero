package utils

import (
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var replacer = stringx.NewReplacer(map[string]string{
	"V": "",
	"v": "",
	"-": ".",
})

// CompareVersions returns true if the first field and the third field are equal, otherwise false.
func CompareVersions(v1, op, v2 string) bool {
	result := compare(v1, v2)
	switch op {
	case "=", "==":
		return result == 0
	case "<":
		return result == -1
	case ">":
		return result == 1
	case "<=":
		return result == -1 || result == 0
	case ">=":
		return result == 0 || result == 1
	}

	return false
}

// return -1 if v1<v2, 0 if they are equal, and 1 if v1>v2
func compare(v1, v2 string) int {
	v1 = replacer.Replace(v1)
	v2 = replacer.Replace(v2)
	fields1 := strings.Split(v1, ".")
	fields2 := strings.Split(v2, ".")
	ver1 := strsToInts(fields1)
	ver2 := strsToInts(fields2)
	shorter := mathx.MinInt(len(ver1), len(ver2))

	for i := 0; i < shorter; i++ {
		if ver1[i] == ver2[i] {
			continue
		} else if ver1[i] < ver2[i] {
			return -1
		} else {
			return 1
		}
	}

	if len(ver1) < len(ver2) {
		return -1
	} else if len(ver1) == len(ver2) {
		return 0
	} else {
		return 1
	}
}

func strsToInts(strs []string) []int64 {
	if len(strs) == 0 {
		return nil
	}

	ret := make([]int64, 0, len(strs))
	for _, str := range strs {
		i, _ := strconv.ParseInt(str, 10, 64)
		ret = append(ret, i)
	}

	return ret
}
