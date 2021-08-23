package template

// Field defines a filed template for types
var Field = `{{.name}} {{.type}} {{.tag}} {{if .hasComment}}// {{.comment}}{{end}}`
