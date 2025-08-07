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
			start, end := val[0], val[len(val)-1]
			if start != '[' && start != '(' {
				return nil, nil, false, false
			}
			if end != ']' && end != ')' {
				return nil, nil, false, false
			}
			exclusiveMinimum = start == '('
			exclusiveMaximum = end == ')'

			content := val[1 : len(val)-1]
			idxColon := strings.Index(content, ":")
			if idxColon < 0 {
				return nil, nil, false, false
			}
			var (
				minStr, maxStr string
				minVal, maxVal *float64
			)
			minStr = util.TrimWhiteSpace(content[:idxColon])
			if len(val) >= idxColon+1 {
				maxStr = util.TrimWhiteSpace(content[idxColon+1:])
			}

			if len(minStr) > 0 {
				min, err := strconv.ParseFloat(minStr, 64)
				if err != nil {
					return nil, nil, false, false
				}
				minVal = &min
			}

			if len(maxStr) > 0 {
				max, err := strconv.ParseFloat(maxStr, 64)
				if err != nil {
					return nil, nil, false, false
				}
				maxVal = &max
			}

			return minVal, maxVal, exclusiveMinimum, exclusiveMaximum
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
			var resp = make([]any, 0)
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

func defValueFromOptions(ctx Context, options []string, apiType spec.Type) any {
	tp := sampleTypeFromGoType(ctx, apiType)
	return valueFromOptions(ctx, options, defFlag, tp)
}

func exampleValueFromOptions(ctx Context, options []string, apiType spec.Type) any {
	tp := sampleTypeFromGoType(ctx, apiType)
	val := valueFromOptions(ctx, options, exampleFlag, tp)
	if val != nil {
		return val
	}
	return defValueFromOptions(ctx, options, apiType)
}

func valueFromOptions(_ Context, options []string, key string, tp string) any {
	if len(options) == 0 {
		return nil
	}
	for _, option := range options {
		if strings.HasPrefix(option, key) {
			s := option[len(key):]
			switch tp {
			case swaggerTypeInteger:
				val, _ := strconv.ParseInt(s, 10, 64)
				return val
			case swaggerTypeBoolean:
				val, _ := strconv.ParseBool(s)
				return val
			case swaggerTypeNumber:
				val, _ := strconv.ParseFloat(s, 64)
				return val
			case swaggerTypeArray:
				return s
			case swaggerTypeString:
				return s
			default:
				return nil
			}
		}
	}
	return nil
}
