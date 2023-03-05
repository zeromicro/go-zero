package cache

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserSchoolModel struct {
	conn sqlc.CachedConn
}

func (m *UserSchoolModel) Do(ctx context.Context) error {
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		withUserModel(session)
	})
}
