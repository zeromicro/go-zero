package redis

import (
	"strings"

	red "github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mapping"
	"github.com/tal-tech/go-zero/core/timex"
)

func process(proc func(red.Cmder) error) func(red.Cmder) error {
	return func(cmd red.Cmder) error {
		start := timex.Now()

		defer func() {
			duration := timex.Since(start)
			if duration > slowThreshold {
				var buf strings.Builder
				for i, arg := range cmd.Args() {
					if i > 0 {
						buf.WriteByte(' ')
					}
					buf.WriteString(mapping.Repr(arg))
				}
				logx.WithDuration(duration).Slowf("[REDIS] slowcall on executing: %s", buf.String())
			}
		}()

		return proc(cmd)
	}
}
