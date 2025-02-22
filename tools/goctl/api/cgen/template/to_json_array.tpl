{{indent $.Indent}}cJSON* v_{{$.VarName}} = cJSON_CreateArray();
{{indent $.Indent}}{
{{indent $.Indent}}    {{$.Items.VarCType}}* v_{{$.VarName}}_items = ({{$.Items.VarCType}}*)({{$.VarExpr}}.items);
{{indent $.Indent}}    for (int i = 0; i < {{$.VarExpr}}.count; ++i) { {{if eq $.Items.VarType "primitive"}}
{{template "to_primitive" $.Items}}{{else if eq $.Items.VarType "object"}}
{{template "to_object" $.Items}}{{else if eq $.Items.VarType "array"}}
{{template "to_array" $.Items}}{{end}}
{{indent $.Indent}}        cJSON_AddItemToArray(v_{{$.VarName}}, v_{{.VarName}}_item);
{{indent $.Indent}}    }
{{indent $.Indent}}}