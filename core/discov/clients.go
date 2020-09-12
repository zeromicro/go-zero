package discov

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/discov/internal"
)

const (
	indexOfKey = iota
	indexOfId
)

const timeToLive int64 = 10

var TimeToLive = timeToLive

func extract(etcdKey string, index int) (string, bool) {
	if index < 0 {
		return "", false
	}

	fields := strings.FieldsFunc(etcdKey, func(ch rune) bool {
		return ch == internal.Delimiter
	})
	if index >= len(fields) {
		return "", false
	}

	return fields[index], true
}

func extractId(etcdKey string) (string, bool) {
	return extract(etcdKey, indexOfId)
}

func makeEtcdKey(key string, id int64) string {
	return fmt.Sprintf("%s%c%d", key, internal.Delimiter, id)
}
