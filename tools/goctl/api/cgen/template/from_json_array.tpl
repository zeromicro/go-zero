{{indent $.Indent}}if (cJSON_IsArray(v_{{$.VarName}})) {
{{indent $.Indent}}    {{$.VarExpr}}.count = cJSON_GetArraySize(v_{{$.VarName}});
{{indent $.Indent}}    {{$.VarExpr}}.items = malloc({{$.VarExpr}}.count * sizeof({{$.Items.VarCType}}));
{{indent $.Indent}}    if ({{$.VarExpr}}.items == NULL) {
{{indent $.Indent}}        goto exit_free;
{{indent $.Indent}}    }
{{indent $.Indent}}    for (int i = 0; i < {{$.VarExpr}}.count; ++i) {
{{indent $.Indent}}        cJSON* v_{{$.VarName}}_item = cJSON_GetArrayItem(v_{{$.VarName}}, i);
{{indent $.Indent}}        {{$.Items.VarCType}}* v_{{$.VarName}}_items = ({{$.Items.VarCType}}*)({{$.VarExpr}}.items);{{if eq $.Items.VarType "primitive"}}
{{template "from_primitive" $.Items}}{{else if eq $.Items.VarType "object"}}
{{template "from_object" $.Items}}{{else if eq $.Items.VarType "array"}}
{{template "from_array" $.Items}}{{end}}
{{indent $.Indent}}    }
{{indent $.Indent}}}