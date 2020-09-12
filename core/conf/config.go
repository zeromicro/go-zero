package conf

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/tal-tech/go-zero/core/mapping"
)

var loaders = map[string]func([]byte, interface{}) error{
	".json": LoadConfigFromJsonBytes,
	".yaml": LoadConfigFromYamlBytes,
	".yml":  LoadConfigFromYamlBytes,
}

func LoadConfig(file string, v interface{}) error {
	if content, err := ioutil.ReadFile(file); err != nil {
		return err
	} else if loader, ok := loaders[path.Ext(file)]; ok {
		return loader(content, v)
	} else {
		return fmt.Errorf("unrecoginized file type: %s", file)
	}
}

func LoadConfigFromJsonBytes(content []byte, v interface{}) error {
	return mapping.UnmarshalJsonBytes(content, v)
}

func LoadConfigFromYamlBytes(content []byte, v interface{}) error {
	return mapping.UnmarshalYamlBytes(content, v)
}

func MustLoad(path string, v interface{}) {
	if err := LoadConfig(path, v); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}
