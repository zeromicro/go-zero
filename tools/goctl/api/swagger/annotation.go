package swagger

import (
	"strconv"

	"github.com/zeromicro/go-zero/tools/goctl/util"
	"google.golang.org/grpc/metadata"
)

func getBoolFromKVOrDefault(properties map[string]string, key string, def bool) bool {
	if len(properties) == 0 {
		return def
	}
	md := metadata.New(properties)
	val := md.Get(key)
	if len(val) == 0 {
		return def
	}
	str := util.Unquote(val[0])
	if len(str) == 0 {
		return def
	}
	res, _ := strconv.ParseBool(str)
	return res
}

func getStringFromKVOrDefault(properties map[string]string, key string, def string) string {
	if len(properties) == 0 {
		return def
	}
	md := metadata.New(properties)
	val := md.Get(key)
	if len(val) == 0 {
		return def
	}
	str := util.Unquote(val[0])
	if len(str) == 0 {
		return def
	}
	return str
}

func getListFromInfoOrDefault(properties map[string]string, key string, def []string) []string {
	if len(properties) == 0 {
		return def
	}
	md := metadata.New(properties)
	val := md.Get(key)
	if len(val) == 0 {
		return def
	}

	str := util.Unquote(val[0])
	if len(str) == 0 {
		return def
	}
	resp := util.FieldsAndTrimSpace(str, commaRune)
	if len(resp) == 0 {
		return def
	}
	return resp
}

func getFirstUsableString(def ...string) string {
	if len(def) == 0 {
		return ""
	}
	for _, val := range def {
		str := util.Unquote(val)
		if len(str) != 0 {
			return str
		}
	}
	return ""
}
