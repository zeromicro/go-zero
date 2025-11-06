package swagger

import (
	"strconv"

	"github.com/zeromicro/go-zero/tools/goctl/util"
	"google.golang.org/grpc/metadata"
)

func getBoolFromKVOrDefault(properties map[string]string, key string, def bool) bool {
	return getOrDefault(properties, key, def, func(str string, def bool) bool {
		res, err := strconv.ParseBool(str)
		if err != nil {
			return def
		}

		return res
	})
}

func getFirstUsableString(def ...string) string {
	if len(def) == 0 {
		return ""
	}

	for _, val := range def {
		// Try to unquote if it's a quoted string
		if str, err := strconv.Unquote(val); err == nil && len(str) != 0 {
			return str
		}

		// Otherwise, use the value as-is if it's not empty
		if len(val) != 0 {
			return val
		}
	}

	return ""
}

func getListFromInfoOrDefault(properties map[string]string, key string, def []string) []string {
	return getOrDefault(properties, key, def, func(str string, def []string) []string {
		resp := util.FieldsAndTrimSpace(str, commaRune)
		if len(resp) == 0 {
			return def
		}
		return resp
	})
}

// getOrDefault abstracts the common logic for fetching, unquoting, and defaulting.
func getOrDefault[T any](properties map[string]string, key string, def T, convert func(string, T) T) T {
	if len(properties) == 0 {
		return def
	}

	md := metadata.New(properties)
	val := md.Get(key)
	if len(val) == 0 {
		return def
	}

	str := val[0]
	if unquoted, err := strconv.Unquote(str); err == nil {
		str = unquoted
	}
	if len(str) == 0 {
		return def
	}

	return convert(str, def)
}

func getStringFromKVOrDefault(properties map[string]string, key string, def string) string {
	return getOrDefault(properties, key, def, func(str string, def string) string {
		return str
	})
}
