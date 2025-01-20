package unigen

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func camelCase(raw string, isPascal bool) string {
	re := regexp.MustCompile("[A-Z_/: -]")
	vs := re.FindAllStringIndex(raw, -1)

	// 全小写
	if len(vs) == 0 {
		return raw
	}

	// 小写开头
	if vs[0][0] > 0 {
		vs = append([][]int{{0, vs[0][0]}}, vs...)
	}

	// 满
	vc := len(vs)
	for i := 0; i < vc; i++ {
		if (i + 1) < len(vs) {
			vs[i][1] = vs[i+1][0]
		} else {
			vs[i][1] = len(raw)
		}
	}

	// 驼峰
	ss := make([]string, len(vs))
	c := cases.Title(language.English)
	for i, v := range vs {
		s := strings.Trim(raw[v[0]:v[1]], "/:_ -")
		if i == 0 && !isPascal {
			ss[i] = strings.ToLower(s)
		} else {
			ss[i] = c.String(s)
		}
	}

	return strings.Join(ss, "")
}
