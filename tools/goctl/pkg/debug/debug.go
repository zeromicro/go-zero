package debug

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/env"
)

func IsDebug() bool {
	goctlDebug := env.GetOr(env.GoctlDebug, "false")
	return strings.EqualFold(goctlDebug, "true")
}
