package util

var headTemplate = `// Code generated by goctl. DO NOT EDIT!
// Source: {{.source}}`

// GetHead returns a code head string with source filename
func GetHead(source string) string {
	buffer, _ := With("head").Parse(headTemplate).Execute(map[string]interface{}{
		"source": source,
	})
	return buffer.String()
}
