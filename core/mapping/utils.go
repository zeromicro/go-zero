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

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/stringx"
)

const (
	defaultOption      = "default"
	envOption          = "env"
	inheritOption      = "inherit"
	stringOption       = "string"
	optionalOption     = "optional"
	optionsOption      = "options"
	rangeOption        = "range"
	optionSeparator    = "|"
	equalToken         = "="
	escapeChar         = '\\'
	leftBracket        = '('
	rightBracket       = ')'
	leftSquareBracket  = '['
	rightSquareBracket = ']'
	segmentSeparator   = ','
	intSize            = 32 << (^uint(0) >> 63) // 32 or 64
)

var (
	errUnsupportedType  = errors.New("unsupported type on setting field value")
	errNumberRange      = errors.New("wrong number range setting")
	errNilSliceElement  = errors.New("null element for slice")
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

// Deref dereferences a type, if pointer type, returns its element type.
func Deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

// Repr returns the string representation of v.
func Repr(v any) string {
	return lang.Repr(v)
}

// SetValue sets target to value, pointers are processed automatically.
func SetValue(tp reflect.Type, value, target reflect.Value) {
	value.Set(convertTypeOfPtr(tp, target))
}

// SetMapIndexValue sets target to value at key position, pointers are processed automatically.
func SetMapIndexValue(tp reflect.Type, value, key, target reflect.Value) {
	value.SetMapIndex(key, convertTypeOfPtr(tp, target))
}

// ValidatePtr validates v if it's a valid pointer.
func ValidatePtr(v reflect.Value) error {
	// sequence is very important, IsNil must be called after checking Kind() with reflect.Ptr,
	// panic otherwise
	if !v.IsValid() || v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("not a valid pointer: %v", v)
	}

	return nil
}

func convertTypeFromString(kind reflect.Kind, str string) (any, error) {
	switch kind {
	case reflect.Bool:
		switch strings.ToLower(str) {
		case "1", "true":
			return true, nil
		case "0", "false":
			return false, nil
		default:
			return false, errTypeMismatch
		}
	case reflect.Int:
		return strconv.ParseInt(str, 10, intSize)
	case reflect.Int8:
		return strconv.ParseInt(str, 10, 8)
	case reflect.Int16:
		return strconv.ParseInt(str, 10, 16)
	case reflect.Int32:
		return strconv.ParseInt(str, 10, 32)
	case reflect.Int64:
		return strconv.ParseInt(str, 10, 64)
	case reflect.Uint:
		return strconv.ParseUint(str, 10, intSize)
	case reflect.Uint8:
		return strconv.ParseUint(str, 10, 8)
	case reflect.Uint16:
		return strconv.ParseUint(str, 10, 16)
	case reflect.Uint32:
		return strconv.ParseUint(str, 10, 32)
	case reflect.Uint64:
		return strconv.ParseUint(str, 10, 64)
	case reflect.Float32:
		return strconv.ParseFloat(str, 32)
	case reflect.Float64:
		return strconv.ParseFloat(str, 64)
	case reflect.String:
		return str, nil
	default:
		return nil, errUnsupportedType
	}
}

func convertTypeOfPtr(tp reflect.Type, target reflect.Value) reflect.Value {
	// keep the original value is a pointer
	if tp.Kind() == reflect.Ptr && target.CanAddr() {
		tp = tp.Elem()
		target = target.Addr()
	}

	for tp.Kind() == reflect.Ptr {
		p := reflect.New(target.Type())
		p.Elem().Set(target)
		target = p
		tp = tp.Elem()
	}

	return target
}

func doParseKeyAndOptions(field reflect.StructField, value string) (string, *fieldOptions, error) {
	segments := parseSegments(value)
	key := strings.TrimSpace(segments[0])
	options := segments[1:]

	if len(options) == 0 {
		return key, nil, nil
	}

	var fieldOpts fieldOptions
	for _, segment := range options {
		option := strings.TrimSpace(segment)
		if err := parseOption(&fieldOpts, field.Name, option); err != nil {
			return "", nil, err
		}
	}

	return key, &fieldOpts, nil
}

// ensureValue ensures nested members not to be nil.
// If pointer value is nil, set to a new value.
func ensureValue(v reflect.Value) reflect.Value {
	for {
		if v.Kind() != reflect.Ptr {
			break
		}

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	return v
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

func isLeftInclude(b byte) (bool, error) {
	switch b {
	case '[':
		return true, nil
	case '(':
		return false, nil
	default:
		return false, errNumberRange
	}
}

func isRightInclude(b byte) (bool, error) {
	switch b {
	case ']':
		return true, nil
	case ')':
		return false, nil
	default:
		return false, errNumberRange
	}
}

func maybeNewValue(fieldType reflect.Type, value reflect.Value) {
	if fieldType.Kind() == reflect.Ptr && value.IsNil() {
		value.Set(reflect.New(value.Type().Elem()))
	}
}

func parseGroupedSegments(val string) []string {
	val = strings.TrimLeftFunc(val, func(r rune) bool {
		return r == leftBracket || r == leftSquareBracket
	})
	val = strings.TrimRightFunc(val, func(r rune) bool {
		return r == rightBracket || r == rightSquareBracket
	})
	return parseSegments(val)
}

// don't modify returned fieldOptions, it's cached and shared among different calls.
func parseKeyAndOptions(tagName string, field reflect.StructField) (string, *fieldOptions, error) {
	value := strings.TrimSpace(field.Tag.Get(tagName))
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

	leftInclude, err := isLeftInclude(str[0])
	if err != nil {
		return nil, err
	}

	str = str[1:]
	if len(str) == 0 {
		return nil, errNumberRange
	}

	rightInclude, err := isRightInclude(str[len(str)-1])
	if err != nil {
		return nil, err
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

	if left > right {
		return nil, errNumberRange
	}

	// [2:2] valid
	// [2:2) invalid
	// (2:2] invalid
	// (2:2) invalid
	if left == right {
		if !leftInclude || !rightInclude {
			return nil, errNumberRange
		}
	}

	return &numberRange{
		left:         left,
		leftInclude:  leftInclude,
		right:        right,
		rightInclude: rightInclude,
	}, nil
}

func parseOption(fieldOpts *fieldOptions, fieldName, option string) error {
	switch {
	case option == inheritOption:
		fieldOpts.Inherit = true
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
			return fmt.Errorf("field %q has wrong optional", fieldName)
		}
	case strings.HasPrefix(option, optionsOption):
		val, err := parseProperty(fieldName, optionsOption, option)
		if err != nil {
			return err
		}

		fieldOpts.Options = parseOptions(val)
	case strings.HasPrefix(option, defaultOption):
		val, err := parseProperty(fieldName, defaultOption, option)
		if err != nil {
			return err
		}

		fieldOpts.Default = val
	case strings.HasPrefix(option, envOption):
		val, err := parseProperty(fieldName, envOption, option)
		if err != nil {
			return err
		}

		fieldOpts.EnvVar = val
	case strings.HasPrefix(option, rangeOption):
		val, err := parseProperty(fieldName, rangeOption, option)
		if err != nil {
			return err
		}

		nr, err := parseNumberRange(val)
		if err != nil {
			return err
		}

		fieldOpts.Range = nr
	}

	return nil
}

// parseOptions parses the given options in tag.
// for example, `json:"name,options=foo|bar"` or `json:"name,options=[foo,bar]"`
func parseOptions(val string) []string {
	if len(val) == 0 {
		return nil
	}

	if val[0] == leftSquareBracket {
		return parseGroupedSegments(val)
	}

	return strings.Split(val, optionSeparator)
}

func parseProperty(field, tag, val string) (string, error) {
	segs := strings.Split(val, equalToken)
	if len(segs) != 2 {
		return "", fmt.Errorf("field %q has wrong tag value %q", field, tag)
	}

	return strings.TrimSpace(segs[1]), nil
}

func parseSegments(val string) []string {
	var segments []string
	var escaped, grouped bool
	var buf strings.Builder

	for _, ch := range val {
		if escaped {
			buf.WriteRune(ch)
			escaped = false
			continue
		}

		switch ch {
		case segmentSeparator:
			if grouped {
				buf.WriteRune(ch)
			} else {
				// need to trim spaces, but we cannot ignore empty string,
				// because the first segment stands for the key might be empty.
				// if ignored, the later tag will be used as the key.
				segments = append(segments, strings.TrimSpace(buf.String()))
				buf.Reset()
			}
		case escapeChar:
			if grouped {
				buf.WriteRune(ch)
			} else {
				escaped = true
			}
		case leftBracket, leftSquareBracket:
			buf.WriteRune(ch)
			grouped = true
		case rightBracket, rightSquareBracket:
			buf.WriteRune(ch)
			grouped = false
		default:
			buf.WriteRune(ch)
		}
	}

	last := strings.TrimSpace(buf.String())
	// ignore last empty string
	if len(last) > 0 {
		segments = append(segments, last)
	}

	return segments
}

func setMatchedPrimitiveValue(kind reflect.Kind, value reflect.Value, v any) error {
	switch kind {
	case reflect.Bool:
		value.SetBool(v.(bool))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(v.(int64))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(v.(uint64))
		return nil
	case reflect.Float32, reflect.Float64:
		value.SetFloat(v.(float64))
		return nil
	case reflect.String:
		value.SetString(v.(string))
		return nil
	default:
		return errUnsupportedType
	}
}

func setValueFromString(kind reflect.Kind, value reflect.Value, str string) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	value = ensureValue(value)
	v, err := convertTypeFromString(kind, str)
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

func toFloat64(v any) (float64, bool) {
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

func validateAndSetValue(kind reflect.Kind, value reflect.Value, str string,
	opts *fieldOptionsWithContext) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	v, err := convertTypeFromString(kind, str)
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

func validateValueInOptions(val any, options []string) error {
	if len(options) > 0 {
		switch v := val.(type) {
		case string:
			if !stringx.Contains(options, v) {
				return fmt.Errorf(`error: value %q is not defined in options "%v"`, v, options)
			}
		default:
			if !stringx.Contains(options, Repr(v)) {
				return fmt.Errorf(`error: value "%v" is not defined in options "%v"`, val, options)
			}
		}
	}

	return nil
}

func validateValueRange(mapValue any, opts *fieldOptionsWithContext) error {
	if opts == nil || opts.Range == nil {
		return nil
	}

	fv, ok := toFloat64(mapValue)
	if !ok {
		return errNumberRange
	}

	return validateNumberRange(fv, opts.Range)
}
