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

func Marshal(val interface{}) (map[string]map[string]interface{}, error) {
	ret := make(map[string]map[string]interface{})
	tp := reflect.TypeOf(val)
	rv := reflect.ValueOf(val)

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		value := rv.Field(i)
		if err := processMember(field, value, ret); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func getTag(field reflect.StructField) (tag string, ok bool) {
	tag, _, ok = strings.Cut(string(field.Tag), tagKVSeparator)
	tag = strings.TrimSpace(tag)
	return
}

func processMember(field reflect.StructField, value reflect.Value,
	collector map[string]map[string]interface{}) error {
	var key string
	var err error
	tag, ok := getTag(field)
	if !ok {
		tag = emptyTag
		key = field.Name
	} else {
		var opt *fieldOptions
		key, opt, err = parseKeyAndOptions(tag, field)
		if err != nil {
			return err
		}

		if err = validate(field, value, opt); err != nil {
			return err
		}
	}

	m, ok := collector[tag]
	if ok {
		m[key] = value.Interface()
	} else {
		m = map[string]interface{}{
			key: value.Interface(),
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

	if len(opt.Options) > 0 {
		if err := validateOptions(value, opt); err != nil {
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
		}
	}
	if !found {
		return fmt.Errorf("field %q not in options", val)
	}

	return nil
}
