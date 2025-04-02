package swagger

import (
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func rangeValueFromOptions(options []string) (minimum *float64, maximum *float64, exclusiveMinimum bool, exclusiveMaximum bool) {
	if len(options) == 0 {
		return nil, nil, false, false
	}
	for _, option := range options {
		if strings.HasPrefix(option, rangeFlag) {
			val := option[6:]
			index := strings.Index(val, ":")
			if index < 0 {
				return nil, nil, false, false
			}
			start, end := val[0], val[len(val)-1]
			exclusiveMinimum = start == '('
			exclusiveMaximum = end == ')'

			content := val[1 : len(val)-1]
			parts := strings.Split(content, ":")
			if len(parts) != 2 {
				return nil, nil, false, false
			}

			minVal, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				return nil, nil, false, false
			}

			maxVal, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, nil, false, false
			}

			return &minVal, &maxVal, exclusiveMinimum, exclusiveMaximum
		}
	}
	return nil, nil, false, false
}

func enumsValueFromOptions(options []string) []any {
	if len(options) == 0 {
		return []any{}
	}
	for _, option := range options {
		if strings.HasPrefix(option, enumFlag) {
			var resp []any
			val := option[8:]
			fields := util.FieldsAndTrimSpace(val, func(r rune) bool {
				return r == '|'
			})
			for _, field := range fields {
				resp = append(resp, field)
			}
			return resp
		}
	}
	return []any{}
}

func defValueFromOptions(options []string, apiType spec.Type) any {
	tp := sampleTypeFromGoType(apiType)
	return valueFromOptions(options, defFlag, tp)
}

func exampleValueFromOptions(options []string, apiType spec.Type) any {
	tp := sampleTypeFromGoType(apiType)
	val := valueFromOptions(options, exampleFlag, tp)
	if val != nil {
		return val
	}
	return defValueFromOptions(options, apiType)
}

func valueFromOptions(options []string, key string, tp string) any {
	if len(options) == 0 {
		return nil
	}
	for _, option := range options {
		if strings.HasPrefix(option, key) {
			s := option[len(key):]
			switch tp {
			case "integer":
				val, _ := strconv.ParseInt(s, 10, 64)
				return val
			case "boolean":
				val, _ := strconv.ParseBool(s)
				return val
			case "number":
				val, _ := strconv.ParseFloat(s, 64)
				return val
			case "string":
				return s
			default:
				return nil
			}
		}
	}
	return nil
}
