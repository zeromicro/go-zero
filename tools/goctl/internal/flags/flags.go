package flags

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

//go:embed default_en.json
var defaultEnFlagConfig []byte

type ConfigLoader struct {
	conf map[string]any
}

func (cl *ConfigLoader) ReadConfig(in io.Reader) error {
	return json.NewDecoder(in).Decode(&cl.conf)
}

func (cl *ConfigLoader) GetString(key string) string {
	keyList := strings.FieldsFunc(key, func(r rune) bool {
		return r == '.'
	})
	var conf = cl.conf
	for idx, k := range keyList {
		val, ok := conf[k]
		if !ok {
			return ""
		}
		if idx < len(keyList)-1 {
			conf, ok = val.(map[string]any)
			if !ok {
				return ""
			}
			continue
		}

		return fmt.Sprint(val)
	}
	return ""
}

type Flags struct {
	loader *ConfigLoader
}

func MustLoad() *Flags {
	loader := &ConfigLoader{
		conf: map[string]any{},
	}
	if err := loader.ReadConfig(bytes.NewBuffer(defaultEnFlagConfig)); err != nil {
		log.Fatal(err)
	}

	return &Flags{
		loader: loader,
	}
}

func setTestData(t *testing.T, data []byte) {
	origin := defaultEnFlagConfig
	defaultEnFlagConfig = data
	t.Cleanup(func() {
		defaultEnFlagConfig = origin
	})
}

func (f *Flags) Get(key string) (string, error) {
	value := f.loader.GetString(key)
	for util.IsTemplateVariable(value) {
		value = util.TemplateVariable(value)
		if value == key {
			return "", fmt.Errorf("the variable can not be self: %q", key)
		}
		return f.Get(value)
	}
	return value, nil
}

var flags *Flags

func Get(key string) string {
	if flags == nil {
		flags = MustLoad()
	}

	v, err := flags.Get(key)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return v
}
