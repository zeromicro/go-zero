package mon

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

type warpSession struct {
	mongo.Session
}

func (w *warpSession) AbortTransaction(ctx context.Context) error {
	ctx, span := startSpan(ctx)
	defer span.End()

	return w.Session.AbortTransaction(ctx)
}

func (w *warpSession) CommitTransaction(ctx context.Context) error {
	ctx, span := startSpan(ctx)
	defer span.End()

	return w.Session.CommitTransaction(ctx)
}

func (w *warpSession) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error), opts ...*mopt.TransactionOptions) (interface{}, error) {
	ctx, span := startSpan(ctx)
	defer span.End()

	return w.Session.WithTransaction(ctx, fn, opts...)
}

func (w *warpSession) EndSession(ctx context.Context) {
	ctx, span := startSpan(ctx)
	defer span.End()

	w.Session.EndSession(ctx)
}
