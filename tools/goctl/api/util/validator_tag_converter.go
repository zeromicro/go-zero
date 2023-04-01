package util

import "strings"

// ConvertValidateTagToSwagger converts the validator tag to swagger comments.
func ConvertValidateTagToSwagger(tagData string) string {
	if tagData == "" || !strings.Contains(tagData, "validate") {
		return ""
	}

	strings.Index(tagData, "validate")
}
