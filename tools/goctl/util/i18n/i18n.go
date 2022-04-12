package i18n

import (
	"os"
	"strings"
)

// IsCN checks whether the user's system language is zh_CN
var IsCN = getSystemLang() == "zh_CN"

func getSystemLang() string {
	lang, ok := os.LookupEnv("LC_ALL")
	if !ok {
		lang, ok = os.LookupEnv("LC_MESSAGES")
	}
	if !ok {
		lang, ok = os.LookupEnv("LANG")
	}
	if !ok {
		return "default"
	}
	ss := strings.Split(lang, ".")
	if len(ss) != 2 {
		return "default"
	}
	return ss[0]
}
