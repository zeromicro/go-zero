package rescue

import (
	"context"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//
//	defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logx.ErrorStack(p)
	}
}

func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		logc.Errorf(ctx, "%+v\n\n%s", p, debug.Stack())
	}
}
