package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/stringx"
)

const (
	defaultKeyName = "key"
	delimiter      = '.'
)

var (
	errTypeMismatch     = errors.New("type mismatch")
	errValueNotSettable = errors.New("value is not settable")
	keyUnmarshaler      = NewUnmarshaler(defaultKeyName)
	cacheKeys           atomic.Value
	cacheKeysLock       sync.Mutex
	durationType        = reflect.TypeOf(time.Duration(0))
	emptyMap            = map[string]interface{}{}
	emptyValue          = reflect.ValueOf(lang.Placeholder)
)

type (
	Unmarshaler struct {
		key  string
		opts unmarshalOptions
	}

	unmarshalOptions struct {
		fromString bool
	}

	keyCache        map[string][]string
	UnmarshalOption func(*unmarshalOptions)
)

func init() {
	cacheKeys.Store(make(keyCache))
}

func NewUnmarshaler(key string, opts ...UnmarshalOption) *Unmarshaler {
	unmarshaler := Unmarshaler{
		key: key,
	}

	for _, opt := range opts {
		opt(&unmarshaler.opts)
	}

	return &unmarshaler
}

func UnmarshalKey(m map[string]interface{}, v interface{}) error {
	return keyUnmarshaler.Unmarshal(m, v)
}

func (u *Unmarshaler) Unmarshal(m map[string]interface{}, v interface{}) error {
	return u.UnmarshalValuer(MapValuer(m), v)
}

func (u *Unmarshaler) UnmarshalValuer(m Valuer, v interface{}) error {
	return u.unmarshalWithFullName(m, v, "")
}

func (u *Unmarshaler) unmarshalWithFullName(m Valuer, v interface{}, fullName string) error {
	rv := reflect.ValueOf(v)
	if err := ValidatePtr(&rv); err != nil {
		return err
	}

	rte := reflect.TypeOf(v).Elem()
	rve := rv.Elem()
	numFields := rte.NumField()
	for i := 0; i < numFields; i++ {
		field := rte.Field(i)
		if usingDifferentKeys(u.key, field) {
			continue
		}

		if err := u.processField(field, rve.Field(i), m, fullName); err != nil {
			return err
		}
	}

	return nil
}

func (u *Unmarshaler) processAnonymousField(field reflect.StructField, value reflect.Value,
	m Valuer, fullName string) error {
	key, options, err := u.parseOptionsWithContext(field, m, fullName)
	if err != nil {
		return err
	}

	if _, hasValue := getValue(m, key); hasValue {
		return fmt.Errorf("fields of %s can't be wrapped inside, because it's anonymous", key)
	}

	if options.optional() {
		return u.processAnonymousFieldOptional(field, value, key, m, fullName)
	} else {
		return u.processAnonymousFieldRequired(field, value, m, fullName)
	}
}

func (u *Unmarshaler) processAnonymousFieldOptional(field reflect.StructField, value reflect.Value,
	key string, m Valuer, fullName string) error {
	var filled bool
	var required int
	var requiredFilled int
	var indirectValue reflect.Value
	fieldType := Deref(field.Type)

	for i := 0; i < fieldType.NumField(); i++ {
		subField := fieldType.Field(i)
		fieldKey, fieldOpts, err := u.parseOptionsWithContext(subField, m, fullName)
		if err != nil {
			return err
		}

		_, hasValue := getValue(m, fieldKey)
		if hasValue {
			if !filled {
				filled = true
				maybeNewValue(field, value)
				indirectValue = reflect.Indirect(value)

			}
			if err = u.processField(subField, indirectValue.Field(i), m, fullName); err != nil {
				return err
			}
		}
		if !fieldOpts.optional() {
			required++
			if hasValue {
				requiredFilled++
			}
		}
	}

	if filled && required != requiredFilled {
		return fmt.Errorf("%s is not fully set", key)
	}

	return nil
}

func (u *Unmarshaler) processAnonymousFieldRequired(field reflect.StructField, value reflect.Value,
	m Valuer, fullName string) error {
	maybeNewValue(field, value)
	fieldType := Deref(field.Type)
	indirectValue := reflect.Indirect(value)

	for i := 0; i < fieldType.NumField(); i++ {
		if err := u.processField(fieldType.Field(i), indirectValue.Field(i), m, fullName); err != nil {
			return err
		}
	}

	return nil
}

func (u *Unmarshaler) processField(field reflect.StructField, value reflect.Value, m Valuer,
	fullName string) error {
	if usingDifferentKeys(u.key, field) {
		return nil
	}

	if field.Anonymous {
		return u.processAnonymousField(field, value, m, fullName)
	} else {
		return u.processNamedField(field, value, m, fullName)
	}
}

func (u *Unmarshaler) processFieldNotFromString(field reflect.StructField, value reflect.Value,
	mapValue interface{}, opts *fieldOptionsWithContext, fullName string) error {
	fieldType := field.Type
	derefedFieldType := Deref(fieldType)
	typeKind := derefedFieldType.Kind()
	valueKind := reflect.TypeOf(mapValue).Kind()

	switch {
	case valueKind == reflect.Map && typeKind == reflect.Struct:
		return u.processFieldStruct(field, value, mapValue, fullName)
	case valueKind == reflect.String && typeKind == reflect.Slice:
		return u.fillSliceFromString(fieldType, value, mapValue, fullName)
	case valueKind == reflect.String && derefedFieldType == durationType:
		return fillDurationValue(fieldType.Kind(), value, mapValue.(string))
	default:
		return u.processFieldPrimitive(field, value, mapValue, opts, fullName)
	}
}

func (u *Unmarshaler) processFieldPrimitive(field reflect.StructField, value reflect.Value,
	mapValue interface{}, opts *fieldOptionsWithContext, fullName string) error {
	fieldType := field.Type
	typeKind := Deref(fieldType).Kind()
	valueKind := reflect.TypeOf(mapValue).Kind()

	switch {
	case typeKind == reflect.Slice && valueKind == reflect.Slice:
		return u.fillSlice(fieldType, value, mapValue)
	case typeKind == reflect.Map && valueKind == reflect.Map:
		return u.fillMap(field, value, mapValue)
	default:
		switch v := mapValue.(type) {
		case json.Number:
			return u.processFieldPrimitiveWithJsonNumber(field, value, v, opts, fullName)
		default:
			if typeKind == valueKind {
				if err := validateValueInOptions(opts.options(), mapValue); err != nil {
					return err
				}

				return fillWithSameType(field, value, mapValue, opts)
			}
		}
	}

	return newTypeMismatchError(fullName)
}

func (u *Unmarshaler) processFieldPrimitiveWithJsonNumber(field reflect.StructField, value reflect.Value,
	v json.Number, opts *fieldOptionsWithContext, fullName string) error {
	fieldType := field.Type
	fieldKind := fieldType.Kind()
	typeKind := Deref(fieldType).Kind()

	if err := validateJsonNumberRange(v, opts); err != nil {
		return err
	}

	if fieldKind == reflect.Ptr {
		value = value.Elem()
	}

	switch typeKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iValue, err := v.Int64()
		if err != nil {
			return err
		}

		value.SetInt(iValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		iValue, err := v.Int64()
		if err != nil {
			return err
		}

		value.SetUint(uint64(iValue))
	case reflect.Float32, reflect.Float64:
		fValue, err := v.Float64()
		if err != nil {
			return err
		}

		value.SetFloat(fValue)
	default:
		return newTypeMismatchError(fullName)
	}

	return nil
}

func (u *Unmarshaler) processFieldStruct(field reflect.StructField, value reflect.Value,
	mapValue interface{}, fullName string) error {
	convertedValue, ok := mapValue.(map[string]interface{})
	if !ok {
		valueKind := reflect.TypeOf(mapValue).Kind()
		return fmt.Errorf("error: field: %s, expect map[string]interface{}, actual %v", fullName, valueKind)
	}

	return u.processFieldStructWithMap(field, value, MapValuer(convertedValue), fullName)
}

func (u *Unmarshaler) processFieldStructWithMap(field reflect.StructField, value reflect.Value,
	m Valuer, fullName string) error {
	if field.Type.Kind() == reflect.Ptr {
		baseType := Deref(field.Type)
		target := reflect.New(baseType).Elem()
		if err := u.unmarshalWithFullName(m, target.Addr().Interface(), fullName); err != nil {
			return err
		}

		value.Set(target.Addr())
	} else if err := u.unmarshalWithFullName(m, value.Addr().Interface(), fullName); err != nil {
		return err
	}

	return nil
}

func (u *Unmarshaler) processNamedField(field reflect.StructField, value reflect.Value,
	m Valuer, fullName string) error {
	key, opts, err := u.parseOptionsWithContext(field, m, fullName)
	if err != nil {
		return err
	}

	fullName = join(fullName, key)
	mapValue, hasValue := getValue(m, key)
	if hasValue {
		return u.processNamedFieldWithValue(field, value, mapValue, key, opts, fullName)
	} else {
		return u.processNamedFieldWithoutValue(field, value, opts, fullName)
	}
}

func (u *Unmarshaler) processNamedFieldWithValue(field reflect.StructField, value reflect.Value,
	mapValue interface{}, key string, opts *fieldOptionsWithContext, fullName string) error {
	if mapValue == nil {
		if opts.optional() {
			return nil
		} else {
			return fmt.Errorf("field %s mustn't be nil", key)
		}
	}

	maybeNewValue(field, value)

	fieldKind := Deref(field.Type).Kind()
	switch fieldKind {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		return u.processFieldNotFromString(field, value, mapValue, opts, fullName)
	default:
		if u.opts.fromString || opts.fromString() {
			valueKind := reflect.TypeOf(mapValue).Kind()
			if valueKind != reflect.String {
				return fmt.Errorf("error: the value in map is not string, but %s", valueKind)
			}

			options := opts.options()
			if len(options) > 0 {
				if !stringx.Contains(options, mapValue.(string)) {
					return fmt.Errorf(`error: value "%s" for field "%s" is not defined in opts "%v"`,
						mapValue, key, options)
				}
			}

			return fillPrimitive(field.Type, value, mapValue, opts, fullName)
		}

		return u.processFieldNotFromString(field, value, mapValue, opts, fullName)
	}
}

func (u *Unmarshaler) processNamedFieldWithoutValue(field reflect.StructField, value reflect.Value,
	opts *fieldOptionsWithContext, fullName string) error {
	derefedType := Deref(field.Type)
	fieldKind := derefedType.Kind()
	if defaultValue, ok := opts.getDefault(); ok {
		if field.Type.Kind() == reflect.Ptr {
			maybeNewValue(field, value)
			value = value.Elem()
		}
		if derefedType == durationType {
			return fillDurationValue(fieldKind, value, defaultValue)
		}
		return setValue(fieldKind, value, defaultValue)
	}

	switch fieldKind {
	case reflect.Array, reflect.Map, reflect.Slice:
		if !opts.optional() {
			return u.processFieldNotFromString(field, value, emptyMap, opts, fullName)
		}
	case reflect.Struct:
		if !opts.optional() {
			required, err := structValueRequired(u.key, derefedType)
			if err != nil {
				return err
			}
			if required {
				return fmt.Errorf("%q is not set", fullName)
			}
			return u.processFieldNotFromString(field, value, emptyMap, opts, fullName)
		}
	default:
		if !opts.optional() {
			return newInitError(fullName)
		}
	}

	return nil
}

func (u *Unmarshaler) fillMap(field reflect.StructField, value reflect.Value, mapValue interface{}) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	fieldKeyType := field.Type.Key()
	fieldElemType := field.Type.Elem()
	targetValue, err := u.generateMap(fieldKeyType, fieldElemType, mapValue)
	if err != nil {
		return err
	}

	value.Set(targetValue)
	return nil
}

func (u *Unmarshaler) fillSlice(fieldType reflect.Type, value reflect.Value, mapValue interface{}) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	baseType := fieldType.Elem()
	baseKind := baseType.Kind()
	dereffedBaseType := Deref(baseType)
	dereffedBaseKind := dereffedBaseType.Kind()
	refValue := reflect.ValueOf(mapValue)
	conv := reflect.MakeSlice(reflect.SliceOf(baseType), refValue.Len(), refValue.Cap())

	var valid bool
	for i := 0; i < refValue.Len(); i++ {
		ithValue := refValue.Index(i).Interface()
		if ithValue == nil {
			continue
		}

		valid = true
		switch dereffedBaseKind {
		case reflect.Struct:
			target := reflect.New(dereffedBaseType)
			if err := u.Unmarshal(ithValue.(map[string]interface{}), target.Interface()); err != nil {
				return err
			}

			if baseKind == reflect.Ptr {
				conv.Index(i).Set(target)
			} else {
				conv.Index(i).Set(target.Elem())
			}
		default:
			if err := u.fillSliceValue(conv, i, dereffedBaseKind, ithValue); err != nil {
				return err
			}
		}
	}

	if valid {
		value.Set(conv)
	}

	return nil
}

func (u *Unmarshaler) fillSliceFromString(fieldType reflect.Type, value reflect.Value,
	mapValue interface{}, fullName string) error {
	var slice []interface{}
	if err := jsonx.UnmarshalFromString(mapValue.(string), &slice); err != nil {
		return err
	}

	baseFieldType := Deref(fieldType.Elem())
	baseFieldKind := baseFieldType.Kind()
	conv := reflect.MakeSlice(reflect.SliceOf(baseFieldType), len(slice), cap(slice))

	for i := 0; i < len(slice); i++ {
		if err := u.fillSliceValue(conv, i, baseFieldKind, slice[i]); err != nil {
			return err
		}
	}

	value.Set(conv)
	return nil
}

func (u *Unmarshaler) fillSliceValue(slice reflect.Value, index int, baseKind reflect.Kind, value interface{}) error {
	switch v := value.(type) {
	case json.Number:
		return setValue(baseKind, slice.Index(index), v.String())
	default:
		// don't need to consider the difference between int, int8, int16, int32, int64,
		// uint, uint8, uint16, uint32, uint64, because they're handled as json.Number.
		if slice.Index(index).Kind() != reflect.TypeOf(value).Kind() {
			return errTypeMismatch
		}

		slice.Index(index).Set(reflect.ValueOf(value))
		return nil
	}
}

func (u *Unmarshaler) generateMap(keyType, elemType reflect.Type, mapValue interface{}) (reflect.Value, error) {
	mapType := reflect.MapOf(keyType, elemType)
	valueType := reflect.TypeOf(mapValue)
	if mapType == valueType {
		return reflect.ValueOf(mapValue), nil
	}

	refValue := reflect.ValueOf(mapValue)
	targetValue := reflect.MakeMapWithSize(mapType, refValue.Len())
	fieldElemKind := elemType.Kind()
	dereffedElemType := Deref(elemType)
	dereffedElemKind := dereffedElemType.Kind()

	for _, key := range refValue.MapKeys() {
		keythValue := refValue.MapIndex(key)
		keythData := keythValue.Interface()

		switch dereffedElemKind {
		case reflect.Slice:
			target := reflect.New(dereffedElemType)
			if err := u.fillSlice(elemType, target.Elem(), keythData); err != nil {
				return emptyValue, err
			}

			targetValue.SetMapIndex(key, target.Elem())
		case reflect.Struct:
			keythMap, ok := keythData.(map[string]interface{})
			if !ok {
				return emptyValue, errTypeMismatch
			}

			target := reflect.New(dereffedElemType)
			if err := u.Unmarshal(keythMap, target.Interface()); err != nil {
				return emptyValue, err
			}

			if fieldElemKind == reflect.Ptr {
				targetValue.SetMapIndex(key, target)
			} else {
				targetValue.SetMapIndex(key, target.Elem())
			}
		case reflect.Map:
			keythMap, ok := keythData.(map[string]interface{})
			if !ok {
				return emptyValue, errTypeMismatch
			}

			innerValue, err := u.generateMap(elemType.Key(), elemType.Elem(), keythMap)
			if err != nil {
				return emptyValue, err
			}

			targetValue.SetMapIndex(key, innerValue)
		default:
			switch v := keythData.(type) {
			case string:
				targetValue.SetMapIndex(key, reflect.ValueOf(v))
			case json.Number:
				target := reflect.New(dereffedElemType)
				if err := setValue(dereffedElemKind, target.Elem(), v.String()); err != nil {
					return emptyValue, err
				}

				targetValue.SetMapIndex(key, target.Elem())
			default:
				targetValue.SetMapIndex(key, keythValue)
			}
		}
	}

	return targetValue, nil
}

func (u *Unmarshaler) parseOptionsWithContext(field reflect.StructField, m Valuer, fullName string) (
	string, *fieldOptionsWithContext, error) {
	key, options, err := parseKeyAndOptions(u.key, field)
	if err != nil {
		return "", nil, err
	} else if options == nil {
		return key, nil, nil
	}

	optsWithContext, err := options.toOptionsWithContext(key, m, fullName)
	if err != nil {
		return "", nil, err
	}

	return key, optsWithContext, nil
}

func WithStringValues() UnmarshalOption {
	return func(opt *unmarshalOptions) {
		opt.fromString = true
	}
}

func fillDurationValue(fieldKind reflect.Kind, value reflect.Value, dur string) error {
	d, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}

	if fieldKind == reflect.Ptr {
		value.Elem().Set(reflect.ValueOf(d))
	} else {
		value.Set(reflect.ValueOf(d))
	}

	return nil
}

func fillPrimitive(fieldType reflect.Type, value reflect.Value, mapValue interface{},
	opts *fieldOptionsWithContext, fullName string) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	baseType := Deref(fieldType)
	if fieldType.Kind() == reflect.Ptr {
		target := reflect.New(baseType).Elem()
		switch mapValue.(type) {
		case string, json.Number:
			value.Set(target.Addr())
			value = target
		}
	}

	switch v := mapValue.(type) {
	case string:
		return validateAndSetValue(baseType.Kind(), value, v, opts)
	case json.Number:
		if err := validateJsonNumberRange(v, opts); err != nil {
			return err
		}
		return setValue(baseType.Kind(), value, v.String())
	default:
		return newTypeMismatchError(fullName)
	}
}

func fillWithSameType(field reflect.StructField, value reflect.Value, mapValue interface{},
	opts *fieldOptionsWithContext) error {
	if !value.CanSet() {
		return errValueNotSettable
	}

	if err := validateValueRange(mapValue, opts); err != nil {
		return err
	}

	if field.Type.Kind() == reflect.Ptr {
		baseType := Deref(field.Type)
		target := reflect.New(baseType).Elem()
		target.Set(reflect.ValueOf(mapValue))
		value.Set(target.Addr())
	} else {
		value.Set(reflect.ValueOf(mapValue))
	}

	return nil
}

// getValue gets the value for the specific key, the key can be in the format of parentKey.childKey
func getValue(m Valuer, key string) (interface{}, bool) {
	keys := readKeys(key)
	return getValueWithChainedKeys(m, keys)
}

func getValueWithChainedKeys(m Valuer, keys []string) (interface{}, bool) {
	if len(keys) == 1 {
		v, ok := m.Value(keys[0])
		return v, ok
	} else if len(keys) > 1 {
		if v, ok := m.Value(keys[0]); ok {
			if nextm, ok := v.(map[string]interface{}); ok {
				return getValueWithChainedKeys(MapValuer(nextm), keys[1:])
			}
		}
	}

	return nil, false
}

func insertKeys(key string, cache []string) {
	cacheKeysLock.Lock()
	defer cacheKeysLock.Unlock()

	keys := cacheKeys.Load().(keyCache)
	// copy the contents into the new map, to guarantee the old map is immutable
	newKeys := make(keyCache)
	for k, v := range keys {
		newKeys[k] = v
	}
	newKeys[key] = cache
	cacheKeys.Store(newKeys)
}

func join(elem ...string) string {
	var builder strings.Builder

	var fillSep bool
	for _, e := range elem {
		if len(e) == 0 {
			continue
		}

		if fillSep {
			builder.WriteByte(delimiter)
		} else {
			fillSep = true
		}

		builder.WriteString(e)
	}

	return builder.String()
}

func newInitError(name string) error {
	return fmt.Errorf("field %s is not set", name)
}

func newTypeMismatchError(name string) error {
	return fmt.Errorf("error: type mismatch for field %s", name)
}

func readKeys(key string) []string {
	cache := cacheKeys.Load().(keyCache)
	if keys, ok := cache[key]; ok {
		return keys
	}

	keys := strings.FieldsFunc(key, func(c rune) bool {
		return c == delimiter
	})
	insertKeys(key, keys)

	return keys
}
