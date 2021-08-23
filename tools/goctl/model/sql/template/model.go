package template

// Model defines a template for model
var Model = `package {{.pkg}}
{{.imports}}
{{.vars}}
{{.types}}
{{.new}}
{{.insert}}
{{.find}}
{{.update}}
{{.delete}}
{{.extraMethod}}
`
