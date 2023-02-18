package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

var _ MultipleSqlConn = (*multipleSqlConn)(nil)

type (
	DBConf struct {
		Leader    string
		Followers []string `json:",optional"`
		// BackSource back to source when all slave data is not available.
		BackSource bool `json:",optional"`
	}

	MultipleSqlConn interface {
		SqlConn
		// LeaderDB returns a leader db.
		LeaderDB() (conn SqlConn, err error)
		// FollowerDB returns a follower db.
		FollowerDB() (conn SqlConn, err error)
	}

	multipleSqlConn struct {
		leader         SqlConn
		enableFollower bool
		p2cPicker      picker
		conf           DBConf
	}
)

// MustNewMultipleSqlConn returns a MultipleSqlConn.
func MustNewMultipleSqlConn(driverName string, conf DBConf) MultipleSqlConn {
	conn, err := NewMultipleSqlConn(driverName, conf)
	logx.Must(err)
	return conn
}

// NewMultipleSqlConn returns a MultipleSqlConn.
func NewMultipleSqlConn(driverName string, conf DBConf) (MultipleSqlConn, error) {
	if err := conf.verify(); err != nil {
		return nil, err
	}

	leader := NewSqlConn(driverName, conf.Leader)
	if err := pingDB(leader); err != nil {
		return nil, fmt.Errorf("%s:%w", conf.Leader, err)
	}

	followers := make([]SqlConn, 0, len(conf.Followers))
	for _, datasource := range conf.Followers {
		follower := NewSqlConn(driverName, datasource)
		if err := pingDB(follower); err != nil {
			return nil, fmt.Errorf("%s:%w", conf.Leader, err)
		}

		followers = append(followers, follower)
	}

	return &multipleSqlConn{
		leader:         leader,
		enableFollower: len(followers) != 0,
		p2cPicker:      newP2cPicker(followers),
	}, nil
}

func (m *multipleSqlConn) Exec(query string, args ...any) (sql.Result, error) {
	return m.ExecCtx(context.Background(), query, args...)
}

func (m *multipleSqlConn) ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.leader.ExecCtx(ctx, query, args...)
}

func (m *multipleSqlConn) Prepare(query string) (StmtSession, error) {
	return m.PrepareCtx(context.Background(), query)
}

func (m *multipleSqlConn) PrepareCtx(ctx context.Context, query string) (StmtSession, error) {
	return m.leader.PrepareCtx(ctx, query)
}

func (m *multipleSqlConn) QueryRow(v any, query string, args ...any) error {
	return m.QueryRowCtx(context.Background(), v, query, args...)
}

func (m *multipleSqlConn) QueryRowCtx(ctx context.Context, v any, query string, args ...any) error {
	db := m.getQueryDb(query)
	return db.query(func(conn SqlConn) error {
		return conn.QueryRowCtx(ctx, v, query, args...)
	})
}

func (m *multipleSqlConn) QueryRowPartial(v any, query string, args ...any) error {
	return m.QueryRowPartialCtx(context.Background(), v, query, args...)
}

func (m *multipleSqlConn) QueryRowPartialCtx(ctx context.Context, v any, query string, args ...any) error {
	db := m.getQueryDb(query)
	return db.query(func(conn SqlConn) error {
		return conn.QueryRowPartialCtx(ctx, v, query, args...)
	})
}

func (m *multipleSqlConn) QueryRows(v any, query string, args ...any) error {
	return m.QueryRowsCtx(context.Background(), v, query, args...)
}

func (m *multipleSqlConn) QueryRowsCtx(ctx context.Context, v any, query string, args ...any) error {
	db := m.getQueryDb(query)
	return db.query(func(conn SqlConn) error {
		return conn.QueryRowsCtx(ctx, v, query, args...)
	})
}

func (m *multipleSqlConn) QueryRowsPartial(v any, query string, args ...any) error {
	return m.QueryRowsPartialCtx(context.Background(), v, query, args...)
}

func (m *multipleSqlConn) QueryRowsPartialCtx(ctx context.Context, v any, query string, args ...any) error {
	db := m.getQueryDb(query)
	return db.query(func(conn SqlConn) error {
		return conn.QueryRowsPartialCtx(ctx, v, query, args...)
	})
}

func (m *multipleSqlConn) RawDB() (*sql.DB, error) {
	return m.leader.RawDB()
}

func (m *multipleSqlConn) Transact(fn func(Session) error) error {
	return m.TransactCtx(context.Background(), func(_ context.Context, session Session) error {
		return fn(session)
	})
}

func (m *multipleSqlConn) TransactCtx(ctx context.Context, fn func(context.Context, Session) error) error {
	return m.leader.TransactCtx(ctx, fn)
}

func (m *multipleSqlConn) containSelect(query string) bool {
	query = strings.TrimSpace(query)
	if len(query) >= 6 {
		return strings.EqualFold(query[:6], "select")
	}

	return false
}

func (m *multipleSqlConn) getQueryDb(query string) queryDb {
	if m.containSelect(query) && m.enableFollower {
		result, err := m.p2cPicker.pick()
		if err == nil {
			return queryDb{
				conn: result.conn,
				done: result.done,
			}
		}

		if !m.conf.BackSource {
			return queryDb{
				error: err,
			}
		}
	}

	return queryDb{conn: m.leader}
}

func (m *multipleSqlConn) LeaderDB() (conn SqlConn, err error) {
	return m.leader, nil
}

func (m *multipleSqlConn) FollowerDB() (SqlConn, error) {
	result, err := m.p2cPicker.pick()
	if err != nil {
		return nil, err
	}

	return result.conn, nil
}

// -------------

type queryDb struct {
	conn  SqlConn
	error error
	done  func(err error)
}

func (q *queryDb) query(query func(conn SqlConn) error) (err error) {
	if q.error != nil {
		return q.error
	}
	defer func() {
		if q.done != nil {
			q.done(err)
		}
	}()

	return query(q.conn)
}

// -------------

func (c DBConf) verify() error {
	if c.Leader == "" {
		return errors.New("leader cannot be empty")
	}

	return nil
}

// -------------

func pingDB(conn SqlConn) error {
	db, err := conn.RawDB()
	if err != nil {
		return err
	}

	return db.Ping()
}
