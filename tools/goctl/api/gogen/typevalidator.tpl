func (v {{.type}}) Validate() error {
{{if ne .localVarContents nil}}
    {{.localVarContents}}
{{end}}

{{if ne .fieldNullContents nil}}
    {{.fieldNullContents}}
{{end}}

{{if ne .fieldValidateContents nil}}
    {{.fieldValidateContents}}
{{end}}
    return nil
}