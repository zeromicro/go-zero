{{indent $.Indent}}cJSON* v_{{$.VarName}} = cJSON_CreateObject();
{{indent $.Indent}}{ {{range $k, $v := $.Pairs}}{{if eq $v.VarType "primitive"}}
{{template "to_primitive" $v}}{{else if eq $v.VarType "array"}}
{{template "to_array" $v}}{{else if eq $v.VarType "object"}}
{{template "to_object" $v}}{{end}}
{{indent $v.Indent}}cJSON_AddItemToObject(v_{{$.VarName}}, "{{$k}}", v_{{$v.VarName}});{{end}}
{{indent $.Indent}}}