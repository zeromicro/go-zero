package sqlx

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

/*
BenchmarkQueryRowsPartial measures the per-row struct scanning cost via sqlmock.
The mock overhead is identical for both sides of any before/after comparison.

Run:
	go test -bench=BenchmarkQueryRowsPartial -benchmem -run='^$' ./core/stores/sqlx/
*/

type benchUser struct {
	ID          *uint64 `db:"id"`
	Username    *string `db:"username"`
	Nickname    *string `db:"nickname"`
	Mobile      *string `db:"mobile"`
	Email       *string `db:"email"`
	Avatar      *string `db:"avatar"`
	Password    *string `db:"password"`
	Status      *int64  `db:"status"`
	Level       *int64  `db:"level"`
	Type        *int64  `db:"type"`
	CityID      *uint64 `db:"city_id"`
	ProvinceID  *uint64 `db:"province_id"`
	AreaID      *uint64 `db:"area_id"`
	LoginCount  *int64  `db:"login_count"`
	LastLoginAt *int64  `db:"last_login_at"`
	CreatedAt   *int64  `db:"created_at"`
	UpdatedAt   *int64  `db:"updated_at"`
	DeletedAt   *int64  `db:"deleted_at"`
	RegisterIP  *string `db:"register_ip"`
	LastLoginIP *string `db:"last_login_ip"`
	InviterID   *uint64 `db:"inviter_id"`
	Channel     *string `db:"channel"`
	DeviceType  *int64  `db:"device_type"`
}

func benchColumns() []string {
	return []string{
		"id", "username", "nickname", "mobile", "email", "avatar", "password",
		"status", "level", "type", "city_id", "province_id", "area_id",
		"login_count", "last_login_at", "created_at", "updated_at", "deleted_at",
		"register_ip", "last_login_ip", "inviter_id", "channel", "device_type",
	}
}

func benchRow(id int) []driver.Value {
	s := func(v string) driver.Value { return v }
	i := func(v int64) driver.Value { return v }
	return []driver.Value{
		i(int64(id)), s("user"), s("nick"), s("13800000000"), s("a@b.c"),
		s("http://avatar"), s("pwd"), i(1), i(2), i(3), i(110000), i(310000),
		i(330100), i(42), i(1723000000), i(1723000001), i(1723000002), i(0),
		s("127.0.0.1"), s("127.0.0.2"), i(9), s("app"), i(1),
	}
}

func benchRows(n int) *sqlmock.Rows {
	cols := benchColumns()
	rows := sqlmock.NewRows(cols)
	for i := 1; i <= n; i++ {
		rows.AddRow(benchRow(i)...)
	}
	return rows
}

func runBenchQueryRowsPartial(b *testing.B, nrows int) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	conn := NewSqlConnFromDB(db)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mock.ExpectQuery("SELECT").WillReturnRows(benchRows(nrows))
		var out []benchUser
		if err := conn.QueryRowsPartialCtx(ctx, &out, "SELECT * FROM user"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQueryRowsPartial1(b *testing.B)   { runBenchQueryRowsPartial(b, 1) }
func BenchmarkQueryRowsPartial10(b *testing.B)  { runBenchQueryRowsPartial(b, 10) }
func BenchmarkQueryRowsPartial100(b *testing.B) { runBenchQueryRowsPartial(b, 100) }
