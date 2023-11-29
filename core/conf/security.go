package conf

import (
	"log"
	"os"
	"reflect"

	"github.com/WqyJh/confcrypt"
)

type SecurityConf struct {
	Enable bool   `json:",default=true"`
	Env    string `json:",default=CONFIG_KEY"` // environment variable name stores the encryption key
}

func findSecurityConfInStruct(v interface{}) (SecurityConf, bool) {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		v = reflect.ValueOf(v).Elem().Interface()
	}
	t := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type == reflect.TypeOf(SecurityConf{}) {
			return reflect.ValueOf(v).FieldByIndex(field.Index).Interface().(SecurityConf), true
		}
	}
	return SecurityConf{}, false
}

func SecurityLoad(path string, v interface{}, opts ...Option) error {
	if err := Load(path, v, opts...); err != nil {
		return err
	}
	c, ok := findSecurityConfInStruct(v)
	if ok && c.Enable {
		key := os.Getenv(c.Env)
		decoded, err := confcrypt.Decode(v, key)
		if err != nil {
			return err
		}
		if reflect.TypeOf(v).Kind() == reflect.Ptr {
			reflect.ValueOf(v).Elem().Set(reflect.ValueOf(decoded).Elem())
			return nil
		}
		reflect.ValueOf(v).Set(reflect.ValueOf(decoded))
	}
	return nil
}

func SecurityMustLoad(path string, v interface{}, opts ...Option) {
	if err := SecurityLoad(path, v, opts...); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}
