package conf

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/internal/encoding"
)

const distanceBetweenUpperAndLower = 32

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadFromJsonBytes,
	".toml": LoadFromTomlBytes,
	".yaml": LoadFromYamlBytes,
	".yml":  LoadFromYamlBytes,
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

	return mapping.UnmarshalJsonMap(toCamelCaseKeyMap(m), v, mapping.WithCanonicalKeyFunc(toCamelCase))
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

func toCamelCase(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	var capNext bool
	boundary := true
	for _, v := range s {
		isCap := v >= 'A' && v <= 'Z'
		isLow := v >= 'a' && v <= 'z'
		if boundary && (isCap || isLow) {
			if capNext {
				if isLow {
					v -= distanceBetweenUpperAndLower
				}
			} else {
				if isCap {
					v += distanceBetweenUpperAndLower
				}
			}
			boundary = false
		}
		if isCap || isLow {
			buf.WriteRune(v)
			capNext = false
		} else if v == ' ' || v == '\t' {
			buf.WriteRune(v)
			capNext = false
			boundary = true
		} else if v == '_' {
			capNext = true
			boundary = true
		} else {
			buf.WriteRune(v)
			capNext = true
		}
	}

	return buf.String()
}

func toCamelCaseInterface(v interface{}) interface{} {
	switch vv := v.(type) {
	case map[string]interface{}:
		return toCamelCaseKeyMap(vv)
	case []interface{}:
		var arr []interface{}
		for _, vvv := range vv {
			arr = append(arr, toCamelCaseInterface(vvv))
		}
		return arr
	default:
		return v
	}
}

func toCamelCaseKeyMap(m map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		res[toCamelCase(k)] = toCamelCaseInterface(v)
	}

	return res
}
