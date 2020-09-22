package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/tal-tech/go-zero/core/stringx"
)

const (
	defaultOption   = "default"
	stringOption    = "string"
	optionalOption  = "optional"
	optionsOption   = "options"
	rangeOption     = "range"
	optionSeparator = "|"
	equalToken      = "="
)

var (
	errUnsupportedType  = errors.New("unsupported type on setting field value")
	errNumberRange      = errors.New("wrong number range setting")
	optionsCache        = make(map[string]optionsCacheValue)
	cacheLock           sync.RWMutex
	structRequiredCache = make(map[reflect.Type]requiredCacheValue)
	structCacheLock     sync.RWMutex
)

type (
	optionsCacheValue struct {
		key     string
		options *fieldOptions
		err     error
	}

	requiredCacheValue struct {
		required bool
		err      error
	}
)

func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

func Repr(v interface{}) string {
	if v == nil {
		return ""
	}

	// if func (v *Type) String() string, we can't use Elem()
	switch vt := v.(type) {
	case fmt.Stringer:
		return vt.String()
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	switch vt := val.Interface().(type) {
	case bool:
		return strconv.FormatBool(vt)
	case error:
		return vt.Error()
	case float32:
		return strconv.FormatFloat(float64(vt), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(vt, 'f', -1, 64)
	case fmt.Stringer:
		return vt.String()
	case int:
		return strconv.Itoa(vt)
	case int8:
		return strconv.Itoa(int(vt))
	case int16:
		return strconv.Itoa(int(vt))
	case int32:
		return strconv.Itoa(int(vt))
	case int64:
		return strconv.FormatInt(vt, 10)
	case string:
		return vt
	case uint:
		return strconv.FormatUint(uint64(vt), 10)
	case uint8:
		return strconv.FormatUint(uint64(vt), 10)
	case uint16:
		return strconv.FormatUint(uint64(vt), 10)
	case uint32:
		return strconv.FormatUint(uint64(vt), 10)
	case uint64:
		return strconv.FormatUint(vt, 10)
	case []byte:
		return string(vt)
	default:
		return fmt.Sprint(val.Interface())
	}
}

func ValidatePtr(v *reflect.Value) error {
	// sequence is very important, IsNil must be called after checking Kind() with reflect.Ptr,
	// panic otherwise
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("not a valid pointer: %v", v)
	}

	return nil
}

func convertType(kind reflect.Kind, str string) (interface{}, error) {
	switch kind {
	case reflect.Bool:
		return str == "1" || strings.ToLower(str) == "true", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(str, 10, 64); err != nil {
			return 0, fmt.Errorf("the value %q cannot parsed as int", str)
		} else {
			return intValue, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(str, 10, 64); err != nil {
			return 0, fmt.Errorf("the value %q cannot parsed as uint", str)
		} else {
			return uintValue, nil
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(str, 64); err != nil {
			return 0, fmt.Errorf("the value %q cannot parsed as float", str)
		} else {
			return floatValue, nil
		}
	case reflect.String:
		return str, nil
	default:
		return nil, errUnsupportedType
	}
}

func doParseKeyAndOptions(field reflect.StructField, value string) (string, *fieldOptions, error) {
	segments := strings.Split(value, ",")
	key := strings.TrimSpace(segments[0])
	options := segments[1:]

	if len(options) > 0 {
		var fieldOpts fieldOptions

		for _, segment := range options {
			option := strings.TrimSpace(segment)
			switch {
			case option == stringOption:
				fieldOpts.FromString = true
			case strings.HasPrefix(option, optionalOption):
				segs := strings.Split(option, equalToken)
				switch len(segs) {
				case 1:
					fieldOpts.Optional = true
				case 2:
					fieldOpts.Optional = true
					fieldOpts.OptionalDep = segs[1]
				default:
					return "", nil, fmt.Errorf("field %s has wrong optional", field.Name)
				}
			case option == optionalOption:
				fieldOpts.Optional = true
			case strings.HasPrefix(option, optionsOption):
				segs := strings.Split(option, equalToken)
				if len(segs) != 2 {
					return "", nil, fmt.Errorf("field %s has wrong options", field.Name)
				} else {
					fieldOpts.Options = strings.Split(segs[1], optionSeparator)
				}
			case strings.HasPrefix(option, defaultOption):
				segs := strings.Split(option, equalToken)
				if len(segs) != 2 {
					return "", nil, fmt.Errorf("field %s has wrong default option", field.Name)
				} else {
					fieldOpts.Default = strings.TrimSpace(segs[1])
				}
			case strings.HasPrefix(option, rangeOption):
				segs := strings.Split(option, equalToken)
				if len(segs) != 2 {
					return "", nil, fmt.Errorf("field %s has wrong range", field.Name)
				}
				if nr, err := parseNumberRange(segs[1]); err != nil {
					return "", nil, err
				} else {
					fieldOpts.Range = nr
				}
			}
		}

		return key, &fieldOpts, nil
	}

	return key, nil, nil
}

func implicitValueRequiredStruct(tag string, tp reflect.Type) (bool, error) {
	numFields := tp.NumField()
	for i := 0; i < numFields; i++ {
		childField := tp.Field(i)
		if usingDifferentKeys(tag, childField) {
			return true, nil
		}

		_, opts, err := parseKeyAndOptions(tag, childField)
		if err != nil {
			return false, err
		}

		if opts == nil {
			if childField.Type.Kind() != reflect.Struct {
				return true, nil
			}

			if required, err := implicitValueRequiredStruct(tag, childField.Type); err != nil {
				return false, err
			} else if required {
				return true, nil
			}
		} else if !opts.Optional && len(opts.Default) == 0 {
			return true, nil
		} else if len(opts.OptionalDep) > 0 && opts.OptionalDep[0] == notSymbol {
			return true, nil
		}
	}

	return false, nil
}

func maybeNewValue(field reflect.StructField, value reflect.Value) {
	if field.Type.Kind() == reflect.Ptr && value.IsNil() {
		value.Set(reflect.New(value.Type().Elem()))
	}
}

// don't modify returned fieldOptions, it's cached and shared among different calls.
func parseKeyAndOptions(tagName string, field reflect.StructField) (string, *fieldOptions, error) {
	value := field.Tag.Get(tagName)
	if len(value) == 0 {
		return field.Name, nil, nil
	}

	cacheLock.RLock()
	cache, ok := optionsCache[value]
	cacheLock.RUnlock()
	if ok {
		return stringx.TakeOne(cache.key, field.Name), cache.options, cache.err
	}

	key, options, err := doParseKeyAndOptions(field, value)
	cacheLock.Lock()
	optionsCache[value] = optionsCacheValue{
		key:     key,
		options: options,
		err:     err,
	}
	cacheLock.Unlock()

	return stringx.TakeOne(key, field.Name), options, err
}

// support below notations:
// [:5] (:5] [:5) (:5)
// [1:] [1:) (1:] (1:)
// [1:5] [1:5) (1:5] (1:5)
func parseNumberRange(str string) (*numberRange, error) {
	if len(str) == 0 {
		return nil, errNumberRange
	}

	var leftInclude bool
	switch str[0] {
	case '[':
		leftInclude = true
	case '(':
		leftInclude = false
	default:
		return nil, errNumberRange
	}

	str = str[1:]
	if len(str) == 0 {
		return nil, errNumberRange
	}

	var rightInclude bool
	switch str[len(str)-1] {
	case ']':
		rightInclude = true
	case ')':
		rightInclude = false
	default:
		return nil, errNumberRange
	}

	str = str[:len(str)-1]
	fields := strings.Split(str, ":")
	if len(fields) != 2 {
		return nil, errNumberRange
	}

	if len(fields[0]) == 0 && len(fields[1]) == 0 {
		return nil, errNumberRange
	}

	var left float64
	if len(fields[0]) > 0 {
		var err error
		if left, err = strconv.ParseFloat(fields[0], 64); err != nil {
			return nil, err
		}
	} else {
		left = -math.MaxFloat64
	}

	var right float64
	if len(fields[1]) > 0 {
		var err error
		if right, err = strconv.ParseFloat(fields[1], 64); err != nil {
			return nil, err
		}
	} else {
		right = math.MaxFloat64
	}

	return &numberRange{
		left:         left,
		leftInclude:  leftInclude,
		right:        right,
		rightInclude: rightInclude,
	}, nil
}

func setMatchedPrimitiveValue(kind reflect.Kind, value reflect.Value, v interface{}) error {
	switch kind {
	case reflect.Bool:
		value.SetBool(v.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(v.(int64))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(v.(uint64))
	case reflect.Float32, reflect.Float64:
		value.SetFloat(v.(float64))
	case reflect.String:
		value.SetString(v.(string))
	default:
		return errUnsupportedType
	}

	return nil
}

func setValue(kind reflect.Kind, value reflect.Value, str string) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	v, err := convertType(kind, str)
	if err != nil {
		return err
	}

	return setMatchedPrimitiveValue(kind, value, v)
}

func structValueRequired(tag string, tp reflect.Type) (bool, error) {
	structCacheLock.RLock()
	val, ok := structRequiredCache[tp]
	structCacheLock.RUnlock()
	if ok {
		return val.required, val.err
	}

	required, err := implicitValueRequiredStruct(tag, tp)
	structCacheLock.Lock()
	structRequiredCache[tp] = requiredCacheValue{
		required: required,
		err:      err,
	}
	structCacheLock.Unlock()

	return required, err
}

func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

func usingDifferentKeys(key string, field reflect.StructField) bool {
	if len(field.Tag) > 0 {
		if _, ok := field.Tag.Lookup(key); !ok {
			return true
		}
	}

	return false
}

func validateAndSetValue(kind reflect.Kind, value reflect.Value, str string, opts *fieldOptionsWithContext) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	v, err := convertType(kind, str)
	if err != nil {
		return err
	}

	if err := validateValueRange(v, opts); err != nil {
		return err
	}

	return setMatchedPrimitiveValue(kind, value, v)
}

func validateJsonNumberRange(v json.Number, opts *fieldOptionsWithContext) error {
	if opts == nil || opts.Range == nil {
		return nil
	}

	fv, err := v.Float64()
	if err != nil {
		return err
	}

	return validateNumberRange(fv, opts.Range)
}

func validateNumberRange(fv float64, nr *numberRange) error {
	if nr == nil {
		return nil
	}

	if (nr.leftInclude && fv < nr.left) || (!nr.leftInclude && fv <= nr.left) {
		return errNumberRange
	}

	if (nr.rightInclude && fv > nr.right) || (!nr.rightInclude && fv >= nr.right) {
		return errNumberRange
	}

	return nil
}

func validateValueInOptions(options []string, value interface{}) error {
	if len(options) > 0 {
		switch v := value.(type) {
		case string:
			if !stringx.Contains(options, v) {
				return fmt.Errorf(`error: value "%s" is not defined in options "%v"`, v, options)
			}
		default:
			if !stringx.Contains(options, Repr(v)) {
				return fmt.Errorf(`error: value "%v" is not defined in options "%v"`, value, options)
			}
		}
	}

	return nil
}

func validateValueRange(mapValue interface{}, opts *fieldOptionsWithContext) error {
	if opts == nil || opts.Range == nil {
		return nil
	}

	fv, ok := toFloat64(mapValue)
	if !ok {
		return errNumberRange
	}

	return validateNumberRange(fv, opts.Range)
}
