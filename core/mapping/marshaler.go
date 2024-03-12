package mapping

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	emptyTag       = ""
	tagKVSeparator = ":"
)

// Marshal marshals the given val and returns the map that contains the fields.
// optional=another is not implemented, and it's hard to implement and not commonly used.
func Marshal(val any) (map[string]map[string]any, error) {
	ret := make(map[string]map[string]any)
	tp := reflect.TypeOf(val)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		value := rv.Field(i)
		if err := processMember(field, value, ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func getTag(field reflect.StructField) (string, bool) {
	tag := string(field.Tag)
	if i := strings.Index(tag, tagKVSeparator); i >= 0 {
		return strings.TrimSpace(tag[:i]), true
	}

	return strings.TrimSpace(tag), false
}

func processMember(field reflect.StructField, value reflect.Value,
	collector map[string]map[string]any) error {
	var key string
	var opt *fieldOptions
	var err error
	tag, ok := getTag(field)
	if !ok {
		tag = emptyTag
		key = field.Name
	} else {
		key, opt, err = parseKeyAndOptions(tag, field)
		if err != nil {
			return err
		}

		if err = validate(field, value, opt); err != nil {
			return err
		}
	}

	val := value.Interface()
	if opt != nil && opt.FromString {
		val = fmt.Sprint(val)
	}

	m, ok := collector[tag]
	if ok {
		m[key] = val
	} else {
		m = map[string]any{
			key: val,
		}
	}
	collector[tag] = m

	return nil
}

func validate(field reflect.StructField, value reflect.Value, opt *fieldOptions) error {
	if opt == nil || !opt.Optional {
		if err := validateOptional(field, value); err != nil {
			return err
		}
	}

	if opt == nil {
		return nil
	}

	if opt.Optional && value.IsZero() {
		return nil
	}

	if len(opt.Options) > 0 {
		if err := validateOptions(value, opt); err != nil {
			return err
		}
	}

	if opt.Range != nil {
		if err := validateRange(value, opt); err != nil {
			return err
		}
	}

	return nil
}

func validateOptional(field reflect.StructField, value reflect.Value) error {
	switch field.Type.Kind() {
	case reflect.Ptr:
		if value.IsNil() {
			return fmt.Errorf("field %q is nil", field.Name)
		}
	case reflect.Array, reflect.Slice, reflect.Map:
		if value.IsNil() || value.Len() == 0 {
			return fmt.Errorf("field %q is empty", field.Name)
		}
	}

	return nil
}

func validateOptions(value reflect.Value, opt *fieldOptions) error {
	var found bool
	val := fmt.Sprint(value.Interface())
	for i := range opt.Options {
		if opt.Options[i] == val {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("field %q not in options", val)
	}

	return nil
}

func validateRange(value reflect.Value, opt *fieldOptions) error {
	var val float64
	switch v := value.Interface().(type) {
	case int:
		val = float64(v)
	case int8:
		val = float64(v)
	case int16:
		val = float64(v)
	case int32:
		val = float64(v)
	case int64:
		val = float64(v)
	case uint:
		val = float64(v)
	case uint8:
		val = float64(v)
	case uint16:
		val = float64(v)
	case uint32:
		val = float64(v)
	case uint64:
		val = float64(v)
	case float32:
		val = float64(v)
	case float64:
		val = v
	default:
		return fmt.Errorf("unknown support type for range %q", value.Type().String())
	}

	// validates [left, right], [left, right), (left, right], (left, right)
	if val < opt.Range.left ||
		(!opt.Range.leftInclude && val == opt.Range.left) ||
		val > opt.Range.right ||
		(!opt.Range.rightInclude && val == opt.Range.right) {
		return fmt.Errorf("%v out of range", value.Interface())
	}

	return nil
}
