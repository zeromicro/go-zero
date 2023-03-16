package flags

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"

	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	flagFileName = "cli.json"
	configType   = "json"
)

//go:embed default.json
var defaultFlagConfig []byte
var flagConfigFile string

func init() {
	goctlHome, err := pathx.GetGoctlHome()
	if err != nil {
		return
	}

	flagConfigFile = filepath.Join(goctlHome, flagFileName)
	_ = os.WriteFile(flagConfigFile, defaultFlagConfig, 0666)
}

func setTestConfigFile(t *testing.T, f string) {
	origin := flagConfigFile
	t.Cleanup(func() {
		flagConfigFile = origin
	})
	flagConfigFile = f
}

type Flags struct {
	v *viper.Viper
}

func MustLoad() *Flags {
	var configContent []byte
	if pathx.FileExists(flagConfigFile) {
		configContent, _ = os.ReadFile(flagConfigFile)
	}
	if len(configContent) == 0 {
		configContent = append(configContent, defaultFlagConfig...)
	}

	v := viper.New()
	v.SetConfigType(configType)
	if err := v.ReadConfig(bytes.NewBuffer(configContent)); err != nil {
		log.Fatal(err)
	}

	return &Flags{
		v: v,
	}
}

func (f *Flags) Get(key string) (string, error) {
	value := f.v.GetString(key)
	for util.IsTemplateVariable(value) {
		value = util.TemplateVariable(value)
		if value == key {
			return "", fmt.Errorf("the variable can not be self: %q", key)
		}
		return f.Get(value)
	}
	return value, nil
}
