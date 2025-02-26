{{indent $.Indent}}if ({{$.VarCheck}}(v_{{$.VarName}}){{if eq $.VarCheck "cJSON_IsString"}} && (NULL != v_{{$.VarName}}->{{$.VarValue}}){{end}}) {
{{indent $.Indent}}    {{$.VarExpr}} = {{if eq $.VarCheck "cJSON_IsString"}}strdup(v_{{$.VarName}}->{{$.VarValue}}){{else}}v_{{$.VarName}}->{{$.VarValue}}{{end}};
{{indent $.Indent}}}