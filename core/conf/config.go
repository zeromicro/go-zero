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

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadFromJsonBytes,
	".toml": LoadFromTomlBytes,
	".yaml": LoadFromYamlBytes,
	".yml":  LoadFromYamlBytes,
}

type fieldInfo struct {
	children map[string]fieldInfo
	mapField *fieldInfo
}

// Load loads config into v from file, .json, .yaml and .yml are acceptable.
func Load(file string, v interface{}, opts ...Option) error {
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
func LoadConfig(file string, v interface{}, opts ...Option) error {
	return Load(file, v, opts...)
}

// LoadFromJsonBytes loads config into v from content json bytes.
func LoadFromJsonBytes(content []byte, v interface{}) error {
	var m map[string]interface{}
	if err := jsonx.Unmarshal(content, &m); err != nil {
		return err
	}

	finfo := buildFieldsInfo(reflect.TypeOf(v))
	lowerCaseKeyMap := toLowerCaseKeyMap(m, finfo)

	return mapping.UnmarshalJsonMap(lowerCaseKeyMap, v, mapping.WithCanonicalKeyFunc(toLowerCase))
}

// LoadConfigFromJsonBytes loads config into v from content json bytes.
// Deprecated: use LoadFromJsonBytes instead.
func LoadConfigFromJsonBytes(content []byte, v interface{}) error {
	return LoadFromJsonBytes(content, v)
}

// LoadFromTomlBytes loads config into v from content toml bytes.
func LoadFromTomlBytes(content []byte, v interface{}) error {
	b, err := encoding.TomlToJson(content)
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(b, v)
}

// LoadFromYamlBytes loads config into v from content yaml bytes.
func LoadFromYamlBytes(content []byte, v interface{}) error {
	b, err := encoding.YamlToJson(content)
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(b, v)
}

// LoadConfigFromYamlBytes loads config into v from content yaml bytes.
// Deprecated: use LoadFromYamlBytes instead.
func LoadConfigFromYamlBytes(content []byte, v interface{}) error {
	return LoadFromYamlBytes(content, v)
}

// MustLoad loads config into v from path, exits on error.
func MustLoad(path string, v interface{}, opts ...Option) {
	if err := Load(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func addOrMergeFields(info fieldInfo, key string, child fieldInfo) {
	if prev, ok := info.children[key]; ok {
		// merge fields
		for k, v := range child.children {
			prev.children[k] = v
		}
		prev.mapField = child.mapField
	} else {
		info.children[key] = child
	}
}

func buildFieldsInfo(tp reflect.Type) fieldInfo {
	tp = mapping.Deref(tp)

	switch tp.Kind() {
	case reflect.Struct:
		return buildStructFieldsInfo(tp)
	case reflect.Array, reflect.Slice:
		return buildFieldsInfo(mapping.Deref(tp.Elem()))
	default:
		return fieldInfo{}
	}
}

func buildStructFieldsInfo(tp reflect.Type) fieldInfo {
	info := fieldInfo{
		children: make(map[string]fieldInfo),
	}

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		name := field.Name
		lowerCaseName := toLowerCase(name)
		ft := mapping.Deref(field.Type)
		// flatten anonymous fields
		if field.Anonymous {
			switch ft.Kind() {
			case reflect.Struct:
				fields := buildFieldsInfo(ft)
				for k, v := range fields.children {
					addOrMergeFields(info, k, v)
				}
				info.mapField = fields.mapField
			case reflect.Map:
				elemField := buildFieldsInfo(mapping.Deref(ft.Elem()))
				info.children[lowerCaseName] = fieldInfo{
					mapField: &elemField,
				}
			default:
				info.children[lowerCaseName] = fieldInfo{
					children: make(map[string]fieldInfo),
				}
			}
			continue
		}

		var finfo fieldInfo
		switch ft.Kind() {
		case reflect.Struct:
			finfo = buildFieldsInfo(ft)
		case reflect.Array, reflect.Slice:
			finfo = buildFieldsInfo(ft.Elem())
		case reflect.Map:
			elemInfo := buildFieldsInfo(mapping.Deref(ft.Elem()))
			finfo.mapField = &elemInfo
		}

		addOrMergeFields(info, lowerCaseName, finfo)
	}

	return info
}

func toLowerCase(s string) string {
	return strings.ToLower(s)
}

func toLowerCaseInterface(v any, info fieldInfo) any {
	switch vv := v.(type) {
	case map[string]interface{}:
		return toLowerCaseKeyMap(vv, info)
	case []interface{}:
		var arr []interface{}
		for _, vvv := range vv {
			arr = append(arr, toLowerCaseInterface(vvv, info))
		}
		return arr
	default:
		return v
	}
}

func toLowerCaseKeyMap(m map[string]any, info fieldInfo) map[string]any {
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
			res[k] = toLowerCaseInterface(v, *info.mapField)
		} else {
			res[k] = v
		}
	}

	return res
}
