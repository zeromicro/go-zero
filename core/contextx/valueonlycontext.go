package contextx

import (
	"context"
	"time"
)

type valueOnlyContext struct {
	context.Context
}

func (valueOnlyContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (valueOnlyContext) Done() <-chan struct{} {
	return nil
}

func (valueOnlyContext) Err() error {
	return nil
}

func ValueOnlyFrom(ctx context.Context) context.Context {
	return valueOnlyContext{
		Context: ctx,
	}
}
