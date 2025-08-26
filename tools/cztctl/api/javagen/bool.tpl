
{{.indent}}{{.decorator}}
{{.indent}}public {{.returnType}} is{{.property}}() {
{{.indent}}	return this.{{.tagValue}};
{{.indent}}}

{{.indent}}public void set{{.property}}({{.type}} {{.propertyValue}}) {
{{.indent}}	this.{{.tagValue}} = {{.propertyValue}};
{{.indent}}}
