package swagger

var (
	tpMapper = map[string]string{
		"uint8":   swaggerTypeInteger,
		"uint16":  swaggerTypeInteger,
		"uint32":  swaggerTypeInteger,
		"uint64":  swaggerTypeInteger,
		"int8":    swaggerTypeInteger,
		"int16":   swaggerTypeInteger,
		"int32":   swaggerTypeInteger,
		"int64":   swaggerTypeInteger,
		"int":     swaggerTypeInteger,
		"uint":    swaggerTypeInteger,
		"byte":    swaggerTypeInteger,
		"float32": swaggerTypeNumber,
		"float64": swaggerTypeNumber,
		"string":  swaggerTypeString,
		"bool":    swaggerTypeBoolean,
	}
	commaRune = func(r rune) bool {
		return r == ','
	}
	slashRune = func(r rune) bool {
		return r == '/'
	}
)
