package gen

import "regexp"

func (g *defaultGenerator) split(source string) []string {
	reg := regexp.MustCompile(createTableFlag)
	index := reg.FindAllStringIndex(source, -1)
	list := make([]string, 0)

	for i := len(index) - 1; i >= 0; i-- {
		subIndex := index[i]
		if len(subIndex) == 0 {
			continue
		}
		start := subIndex[0]
		ddl := source[start:]
		list = append(list, ddl)
		source = source[:start]
	}

	return list
}
