package util

import "github.com/tal-tech/go-zero/tools/goctl/api/spec"

func GetAnnotationValue(annos []spec.Annotation, key, field string) (string, bool) {
	for _, anno := range annos {
		if anno.Name == key {
			value, ok := anno.Properties[field]
			return value, ok
		}
	}
	return "", false
}
