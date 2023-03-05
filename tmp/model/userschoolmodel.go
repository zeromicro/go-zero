package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type TransactionModel struct {
	conn sqlx.SqlConn
}

func (m *TransactionModel) InsertSchoolAndUpdateUser(ctx context.Context, u *User, s *School) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		conn := sqlx.WithSession(session)
		r, err := newUserModel(conn).Insert(ctx, u)
		if err != nil {
			return err
		}
		id, err := r.LastInsertId()
		if err != nil {
			return err
		}
		s.UserId = id

		err = newSchoolModel(conn).Update(ctx, s)
		return err
	})
}
