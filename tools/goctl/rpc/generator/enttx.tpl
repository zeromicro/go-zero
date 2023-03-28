package entx

import (
	"context"
	"fmt"

    "github.com/zeromicro/go-zero/core/logx"

	"{{ .package}}/ent"
)


// WithTx uses transaction in ent.
func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		logx.Errorw("failed to start transaction", logx.Field("detail", err.Error()))
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rollBackErr)
		}
		logx.Errorw("errors occur in transaction", logx.Field("detail", err.Error()))
		return err
	}
	if err := tx.Commit(); err != nil {
		logx.Errorw("failed to commit transaction", logx.Field("detail", err.Error()))
		return err
	}
	return nil
}
