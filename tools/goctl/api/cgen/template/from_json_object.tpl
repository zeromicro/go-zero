{{indent $.Indent}}if (cJSON_IsObject(v_{{$.VarName}})) { {{range $k, $v := $.Pairs}}
{{indent $.Indent}}    cJSON* v_{{$v.VarName}} = cJSON_GetObjectItemCaseSensitive(v_{{$.VarName}}, "{{$k}}");{{if eq $v.VarType "primitive"}}
{{template "from_primitive" $v}}{{else if eq $v.VarType "array"}}
{{template "from_array" $v}}{{else if eq $v.VarType "object"}}
{{template "from_object" $v}}{{end}}{{end}}
{{indent $.Indent}}}