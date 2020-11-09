package template

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
