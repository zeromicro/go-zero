{ {{range $i, $m := .Members}}
    {{Indent $.Indent}}{{$m.PropertyName}}{{$m.OptionalTag}}: {{$m.PropertyType}};{{if $m.Comment}} // {{$m.Comment}}{{end}}{{end}}
{{Indent $.Indent}}}