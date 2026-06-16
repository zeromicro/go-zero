package swagger

import "strings"

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

// extractPathPlaceholders
// Supports both :id style (go-zero format) and {id} style (OpenAPI format).
// For example: "/foo/:id/bar/{namespace}" returns ["id", "namespace"]
func extractPathPlaceholders(path string) map[string]bool {
	placeholders := make(map[string]bool)
	items := strings.FieldsFunc(path, slashRune)

	for _, item := range items {
		// Handle :id style placeholders
		if strings.HasPrefix(item, ":") {
			name := strings.TrimPrefix(item, ":")
			if name != "" {
				placeholders[name] = true
			}
		}

		// Handle {id} style placeholders
		if strings.HasPrefix(item, "{") && strings.HasSuffix(item, "}") {
			name := strings.TrimSuffix(strings.TrimPrefix(item, "{"), "}")
			if name != "" {
				placeholders[name] = true
			}
		}
	}

	return placeholders
}
