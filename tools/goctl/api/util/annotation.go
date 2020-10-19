package util

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

func GetAnnotationValue(annotations []spec.Annotation, key, field string) (string, bool) {
	for _, annotation := range annotations {
		if annotation.Name == field && len(annotation.Value) > 0 {
			return annotation.Value, true
		}
		if annotation.Name == key {
			value, ok := annotation.Properties[field]
			return strings.TrimSpace(value), ok
		}
	}
	return "", false
}
