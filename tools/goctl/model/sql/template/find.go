package sqltemplate

// 通过id查询
var FindOne = `
func (m *{{.upperObject}}Model) FindOne({{.lowerPrimaryKey}} {{.dataType}}) (*{{.upperObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperObject}}
	err := m.QueryRow(&resp, {{.cacheKeyVariable}}, func(conn sqlx.SqlConn, v interface{}) error {
		query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.snakePrimaryKey}} = ? limit 1` + "`" + `
		return conn.QueryRow(v, query, {{.lowerPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.snakePrimaryKey}} = ? limit 1` + "`" + `
	var resp {{.upperObject}}
	err := m.QueryRowNoCache(&resp, query, {{.lowerPrimaryKey}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}
}
`

// 通过指定字段查询
var FindOneByField = `
func (m *{{.upperObject}}Model) FindOneBy{{.upperFields}}({{.in}}) (*{{.upperObject}}, error) {
	{{if .onlyOneFiled}}{{if .withCache}}{{.cacheKey}}
	var resp {{.upperObject}}
	err := m.QueryRowIndex(&resp, {{.cacheKeyVariable}}, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", {{.primaryKeyDefine}}, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.snakeField}} = ? limit 1` + "`" + `
		if err := conn.QueryRow(&resp, query, {{.lowerField}}); err != nil {
			return nil, err
		}
		return resp.{{.UpperPrimaryKey}}, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.primarySnakeField}} = ? limit 1` + "`" + `
		return conn.QueryRow(v, query, primary)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}var resp {{.upperObject}}
	query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.expression}} limit 1` + "`" + `
	err := m.QueryRowNoCache(&resp, query, {{.expressionValues}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}{{else}}var resp {{.upperObject}}
	query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.expression}} limit 1` + "`" + `
	err := m.QueryRowNoCache(&resp, query, {{.expressionValues}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{end}}
}
`

// 查询all
var FindAllByField = `
func (m *{{.upperObject}}Model) FindAllBy{{.upperFields}}({{.in}}) ([]*{{.upperObject}}, error) {
	var resp []*{{.upperObject}}
	query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.expression}}` + "`" + `
	err := m.QueryRowsNoCache(&resp, query, {{.expressionValues}})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
`

// limit分页查询
var FindLimitByField = `
func (m *{{.upperObject}}Model) FindLimitBy{{.upperFields}}({{.in}}, page, limit int) ([]*{{.upperObject}}, error) {
	var resp []*{{.upperObject}}
	query := ` + "`" + `select ` + "`" + ` + {{.lowerObject}}Rows + ` + "`" + `from ` + "` + " + `m.table ` + " + `" + ` where {{.expression}} order by {{.sortExpression}} limit ?,?` + "`" + `
	err := m.QueryRowsNoCache(&resp, query, {{.expressionValues}}, (page-1)*limit, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *{{.upperObject}}Model) FindAllCountBy{{.upperFields}}({{.in}}) (int64, error) {
	var count int64
	query := ` + "`" + `select count(1)  from ` + "` + " + `m.table ` + " + `" + ` where {{.expression}} ` + "`" + `
	err := m.QueryRowNoCache(&count, query, {{.expressionValues}})
	if err != nil {
		return 0, err
	}
	return count, nil
}
`
