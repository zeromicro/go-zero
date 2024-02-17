package conf

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/internal/encoding"
)

const (
	jsonTagKey = "json"
	jsonTagSep = ','
)

var (
	fillDefaultUnmarshaler = mapping.NewUnmarshaler(jsonTagKey, mapping.WithDefault())
	loaders                = map[string]func([]byte, any) error{
		".json": LoadFromJsonBytes,
		".toml": LoadFromTomlBytes,
		".yaml": LoadFromYamlBytes,
		".yml":  LoadFromYamlBytes,
	}
)

// children and mapField should not be both filled.
// named fields and map cannot be bound to the same field name.
type fieldInfo struct {
	children map[string]*fieldInfo
	mapField *fieldInfo
}

// FillDefault fills the default values for the given v,
// and the premise is that the value of v must be guaranteed to be empty.
func FillDefault(v any) error {
	return fillDefaultUnmarshaler.Unmarshal(map[string]any{}, v)
}

// Load loads config into v from file, .json, .yaml and .yml are acceptable.
func Load(file string, v any, opts ...Option) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	loader, ok := loaders[strings.ToLower(path.Ext(file))]
	if !ok {
		return fmt.Errorf("unrecognized file type: %s", file)
	}

	var opt options
	for _, o := range opts {
		o(&opt)
	}

	if opt.env {
		return loader([]byte(os.ExpandEnv(string(content))), v)
	}

	return loader(content, v)
}

// LoadConfig loads config into v from file, .json, .yaml and .yml are acceptable.
// Deprecated: use Load instead.
func LoadConfig(file string, v any, opts ...Option) error {
	return Load(file, v, opts...)
}

// LoadFromJsonBytes loads config into v from content json bytes.
func LoadFromJsonBytes(content []byte, v any) error {
	info, err := buildFieldsInfo(reflect.TypeOf(v), "")
	if err != nil {
		return err
	}

	var m map[string]any
	if err = jsonx.Unmarshal(content, &m); err != nil {
		return err
	}

	lowerCaseKeyMap := toLowerCaseKeyMap(m, info)

	return mapping.UnmarshalJsonMap(lowerCaseKeyMap, v, mapping.WithCanonicalKeyFunc(toLowerCase))
}

// LoadConfigFromJsonBytes loads config into v from content json bytes.
// Deprecated: use LoadFromJsonBytes instead.
func LoadConfigFromJsonBytes(content []byte, v any) error {
	return LoadFromJsonBytes(content, v)
}

// LoadFromTomlBytes loads config into v from content toml bytes.
func LoadFromTomlBytes(content []byte, v any) error {
	b, err := encoding.TomlToJson(content)
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(b, v)
}

// LoadFromYamlBytes loads config into v from content yaml bytes.
func LoadFromYamlBytes(content []byte, v any) error {
	b, err := encoding.YamlToJson(content)
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(b, v)
}

// LoadConfigFromYamlBytes loads config into v from content yaml bytes.
// Deprecated: use LoadFromYamlBytes instead.
func LoadConfigFromYamlBytes(content []byte, v any) error {
	return LoadFromYamlBytes(content, v)
}

// MustLoad loads config into v from path, exits on error.
func MustLoad(path string, v any, opts ...Option) {
	if err := Load(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func addOrMergeFields(info *fieldInfo, key string, child *fieldInfo, fullName string) error {
	if prev, ok := info.children[key]; ok {
		if child.mapField != nil {
			return newConflictKeyError(fullName)
		}

		if err := mergeFields(prev, key, child.children, fullName); err != nil {
			return err
		}
	} else {
		info.children[key] = child
	}

	return nil
}

func buildAnonymousFieldInfo(info *fieldInfo, lowerCaseName string, ft reflect.Type, fullName string) error {
	switch ft.Kind() {
	case reflect.Struct:
		fields, err := buildFieldsInfo(ft, fullName)
		if err != nil {
			return err
		}

		for k, v := range fields.children {
			if err = addOrMergeFields(info, k, v, fullName); err != nil {
				return err
			}
		}
	case reflect.Map:
		elemField, err := buildFieldsInfo(mapping.Deref(ft.Elem()), fullName)
		if err != nil {
			return err
		}

		if _, ok := info.children[lowerCaseName]; ok {
			return newConflictKeyError(fullName)
		}

		info.children[lowerCaseName] = &fieldInfo{
			children: make(map[string]*fieldInfo),
			mapField: elemField,
		}
	default:
		if _, ok := info.children[lowerCaseName]; ok {
			return newConflictKeyError(fullName)
		}

		info.children[lowerCaseName] = &fieldInfo{
			children: make(map[string]*fieldInfo),
		}
	}

	return nil
}

func buildFieldsInfo(tp reflect.Type, fullName string) (*fieldInfo, error) {
	tp = mapping.Deref(tp)

	switch tp.Kind() {
	case reflect.Struct:
		return buildStructFieldsInfo(tp, fullName)
	case reflect.Array, reflect.Slice:
		return buildFieldsInfo(mapping.Deref(tp.Elem()), fullName)
	case reflect.Chan, reflect.Func:
		return nil, fmt.Errorf("unsupported type: %s", tp.Kind())
	default:
		return &fieldInfo{
			children: make(map[string]*fieldInfo),
		}, nil
	}
}

func buildNamedFieldInfo(info *fieldInfo, lowerCaseName string, ft reflect.Type, fullName string) error {
	var finfo *fieldInfo
	var err error

	switch ft.Kind() {
	case reflect.Struct:
		finfo, err = buildFieldsInfo(ft, fullName)
		if err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:
		finfo, err = buildFieldsInfo(ft.Elem(), fullName)
		if err != nil {
			return err
		}
	case reflect.Map:
		elemInfo, err := buildFieldsInfo(mapping.Deref(ft.Elem()), fullName)
		if err != nil {
			return err
		}

		finfo = &fieldInfo{
			children: make(map[string]*fieldInfo),
			mapField: elemInfo,
		}
	default:
		finfo, err = buildFieldsInfo(ft, fullName)
		if err != nil {
			return err
		}
	}

	return addOrMergeFields(info, lowerCaseName, finfo, fullName)
}

func buildStructFieldsInfo(tp reflect.Type, fullName string) (*fieldInfo, error) {
	info := &fieldInfo{
		children: make(map[string]*fieldInfo),
	}

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		if !field.IsExported() {
			continue
		}

		name := getTagName(field)
		lowerCaseName := toLowerCase(name)
		ft := mapping.Deref(field.Type)
		// flatten anonymous fields
		if field.Anonymous {
			if err := buildAnonymousFieldInfo(info, lowerCaseName, ft,
				getFullName(fullName, lowerCaseName)); err != nil {
				return nil, err
			}
		} else if err := buildNamedFieldInfo(info, lowerCaseName, ft,
			getFullName(fullName, lowerCaseName)); err != nil {
			return nil, err
		}
	}

	return info, nil
}

// getTagName get the tag name of the given field, if no tag name, use file.Name.
// field.Name is returned on tags like `json:""` and `json:",optional"`.
func getTagName(field reflect.StructField) string {
	if tag, ok := field.Tag.Lookup(jsonTagKey); ok {
		if pos := strings.IndexByte(tag, jsonTagSep); pos >= 0 {
			tag = tag[:pos]
		}

		tag = strings.TrimSpace(tag)
		if len(tag) > 0 {
			return tag
		}
	}

	return field.Name
}

func mergeFields(prev *fieldInfo, key string, children map[string]*fieldInfo, fullName string) error {
	if len(prev.children) == 0 || len(children) == 0 {
		return newConflictKeyError(fullName)
	}

	// merge fields
	for k, v := range children {
		if _, ok := prev.children[k]; ok {
			return newConflictKeyError(fullName)
		}

		prev.children[k] = v
	}

	return nil
}

func toLowerCase(s string) string {
	return strings.ToLower(s)
}

func toLowerCaseInterface(v any, info *fieldInfo) any {
	switch vv := v.(type) {
	case map[string]any:
		return toLowerCaseKeyMap(vv, info)
	case []any:
		var arr []any
		for _, vvv := range vv {
			arr = append(arr, toLowerCaseInterface(vvv, info))
		}
		return arr
	default:
		return v
	}
}

func toLowerCaseKeyMap(m map[string]any, info *fieldInfo) map[string]any {
	res := make(map[string]any)

	for k, v := range m {
		ti, ok := info.children[k]
		if ok {
			res[k] = toLowerCaseInterface(v, ti)
			continue
		}

		lk := toLowerCase(k)
		if ti, ok = info.children[lk]; ok {
			res[lk] = toLowerCaseInterface(v, ti)
		} else if info.mapField != nil {
			res[k] = toLowerCaseInterface(v, info.mapField)
		} else {
			res[k] = v
		}
	}

	return res
}

type conflictKeyError struct {
	key string
}

func newConflictKeyError(key string) conflictKeyError {
	return conflictKeyError{key: key}
}

func (e conflictKeyError) Error() string {
	return fmt.Sprintf("conflict key %s, pay attention to anonymous fields", e.key)
}

func getFullName(parent, child string) string {
	if len(parent) == 0 {
		return child
	}

	return strings.Join([]string{parent, child}, ".")
}
