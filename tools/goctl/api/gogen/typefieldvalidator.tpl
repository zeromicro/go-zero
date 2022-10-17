{{if eq .type  "number"}}
    {{if and .exclusiveMin .shouldValidMin}}
        if {{.property}} <= {{.min}} {
            return  fmt.Errorf("property: {{.propertyName}}  can not less than min: {{.min}}, actual: %d", {{.property}})
        }
    {{else if .shouldValidMin}}
        if {{.property}} < {{.min}} {
            return fmt.Errorf("property: {{.propertyName}} can not less than min: {{.min}}, actual: %d", {{.property}})
        }
    {{else if .notAllowEmptyValue}}
        if {{.property}} <= 0 {
           return  fmt.Errorf("property: {{.propertyName}}  can not less than min: {{.min}}, actual: %d", {{.property}})
        }
    {{end}}

    {{if and .exclusiveMax .shouldValidMax}}
        if {{.property}} >= {{.max}} {
            return fmt.Errorf("property: {{.propertyName}} can not more than man: {{.max}}, actual: %d", {{.property}})
        }
    {{else if .shouldValidMax}}
        if {{.property}} > {{.max}} {
            return fmt.Errorf("property: {{.propertyName}} can not more than min: {{.max}}, actual: %d", {{.property}})
        }
    {{end}}

    {{if ne .multipleOf nil}}
        // 是multipleOf的整数倍
        if {{.property}} % {{.multipleOf}} != 0 {
           return fmt.Errorf("property: {{.propertyName}} not multiple of: {{.multipleOf}}, actual: %d", {{.property}})
        }
    {{end}}
{{end}}

{{if eq .type  "string"}}
    {{if ne .minLength nil}}
        if len({{.property}}) > {{.minLength}} {
            return fmt.Errorf("property: {{.propertyName}} string length can not less than min:  {{.minLength}}, actual: %s", {{.property}})
        }
    {{else if .notAllowEmptyValue}}
        if len({{.property}}) <= 0 {
           return fmt.Errorf("property: {{.propertyName}} string length can not less than min:  {{.minLength}}, actual: %s", {{.property}})
        }
    {{end}}

    {{if .shouldValidMaxLength}}
        if len({{.property}}) > {{.maxLength}} {
            return fmt.Errorf("property: {{.propertyName}} string length: %s can not more than min: {{.maxLength}}, actual: %s", {{.property}})
        }
    {{end}}

    {{if ne .reg nil}}
	    if {{.reg}}.MatchString({{.property}}) {
		    return fmt.Errorf("property: {{.propertyName}} not match regex: '{{.reg}}', actual: %s", {{.property}})
	    }
    {{end}}
{{end}}

{{if eq .type "array"}}
    {{if ne .minItems nil }}
        if len({{.property}}) < {{.minItems}} {
            return fmt.Errorf("property: {{.propertyName}} array length can not less than min: {{.maxItems}}, actual: %d", len({{.property}}))
        }
    {{else if .notAllowEmptyValue}}
        if len({{.property}}) <= 0 {
           return fmt.Errorf("property:{{.propertyName}} array length can not less than min:  {{.minLength}}, actual: %s", {{.property}})
        }
    {{end}}

    {{if ne .maxItems nil}}
        if len({{.property}}) > {{.maxItems}} {
            return fmt.Errorf("property: {{.propertyName}} array length can not more than min: {{.maxItems}}, actual: %d", len({{.property}}))
        }
    {{end}}

    {{if .shouldValidUniqueItems}}
        arrayRepeatMap = make(map[interface{}]struct{})
        for _, item := range {{.property}}{
            if _, ok := arrayRepeatMap[item]; ok {
        	    return fmt.Errorf("property: {{.propertyName}} item: %v is repeated, array: %v", item, {{.property}})
            }
            arrayRepeatMap[item] = struct{}{}
        }
    {{end}}

    {{if ne .itemValidateContent nil}}
        for _, {{.item}} := range {{.property}} {
            {{.itemValidateContent}}
        }
    {{end}}
{{end}}

{{if eq .type "map"}}
    {{if ne .valueValidateContent nil}}
        for _, {{.value}} := range {{.property}} {
            {{.valueValidateContent}}
        }
    {{end}}
{{end}}

{{if eq .type "object"}}
    if err := {{.property}}.Validate(); err != nil {
        return fmt.Errorf("property: {{.propertyName}} validate err: %w", err)
    }
{{end}}

{{if .enumValidate}}
    enumExist = false
    {{.enumDefine}}
    for _, item := range enumArray {
        if item == {{.property}} {
            enumExist = true
            break
        }
    }
    if !enumExist {
        return fmt.Errorf("property: {{.propertyName}} not in %v", enumArray)
    }
{{end}}