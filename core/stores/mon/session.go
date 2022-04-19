package mon

import (
	"context"

	"github.com/zeromicro/go-zero/core/breaker"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

type warpSession struct {
	mongo.Session
	brk breaker.Breaker
}

func (w *warpSession) AbortTransaction(ctx context.Context) error {
	ctx, span := startSpan(ctx)
	defer span.End()

	return w.brk.DoWithAcceptable(func() error {
		return w.Session.AbortTransaction(ctx)
	}, acceptable)
}

func (w *warpSession) CommitTransaction(ctx context.Context) error {
	ctx, span := startSpan(ctx)
	defer span.End()

	return w.brk.DoWithAcceptable(func() error {
		return w.Session.CommitTransaction(ctx)
	}, acceptable)
}

func (w *warpSession) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error), opts ...*mopt.TransactionOptions) (res interface{}, err error) {
	ctx, span := startSpan(ctx)
	defer span.End()

	err = w.brk.DoWithAcceptable(func() error {
		res, err = w.Session.WithTransaction(ctx, fn, opts...)
		return err
	}, acceptable)

	return
}

func (w *warpSession) EndSession(ctx context.Context) {
	ctx, span := startSpan(ctx)
	defer span.End()
	_ = w.brk.DoWithAcceptable(func() error {
		w.Session.EndSession(ctx)
		return nil
	}, acceptable)
}
