package utils

import (
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

// GO 启动 goroutine 防止崩溃
func GO(ctx context.Context, f func(context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logx.ErrorStack(r)
			}
		}()
		f(ctx)
	}()
}
