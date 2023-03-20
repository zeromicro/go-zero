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
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	flagFileName = "cli-%s.json"
	configType   = "json"
)

//go:embed default_en.json
var defaultEnFlagConfig []byte

//go:embed default_zh.json
var defaultZhFlagConfig []byte

func Init() {
	flagConfigFile := getFlagConfigFile()
	_ = pathx.MkdirIfNotExist(filepath.Dir(flagConfigFile))
	if pathx.FileExists(flagConfigFile) {
		return
	}

	locale := env.GetOr(env.GoctlLocale, env.DefaultLocale)
	var configForLocale []byte
	if locale == env.LocaleSimplifiedChinese {
		configForLocale = append(configForLocale, defaultZhFlagConfig...)
	} else {
		configForLocale = append(configForLocale, defaultEnFlagConfig...)
	}

	_ = os.WriteFile(flagConfigFile, configForLocale, 0666)
}

func getFlagConfigFile() string {
	goctlHome, err := pathx.GetGoctlHome()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	var locale string
	locale = env.Get(env.GoctlLocale)
	if len(locale) == 0 {
		locale = "en"
	}

	return filepath.Join(goctlHome, version.BuildVersion, fmt.Sprintf(flagFileName, locale))
}

var flagConfigFile string

func setTestConfigFile(t *testing.T, f string) {
	origin := getFlagConfigFile()
	t.Cleanup(func() {
		flagConfigFile = origin
	})
	flagConfigFile = f
}

type Flags struct {
	v *viper.Viper
}

func MustLoad() *Flags {
	if len(flagConfigFile) == 0 {
		flagConfigFile = getFlagConfigFile()
	}

	var configContent []byte
	if pathx.FileExists(flagConfigFile) {
		configContent, _ = os.ReadFile(flagConfigFile)
	}

	if len(configContent) == 0 {
		configContent = append(configContent, defaultEnFlagConfig...)
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
