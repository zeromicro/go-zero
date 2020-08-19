package template

// 通过id查询
var FindOne = `
func (m *{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRow(&resp, {{.cacheKeyVariable}}, func(conn sqlx.SqlConn, v interface{}) error {
		query := ` + "`" + `select ` + "`" + ` + {{.lowerStartCamelObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.originalPrimaryKey}} = ? limit 1` + "`" + `
		return conn.QueryRow(v, query, {{.lowerStartCamelPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := ` + "`" + `select ` + "`" + ` + {{.lowerStartCamelObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.originalPrimaryKey}} = ? limit 1` + "`" + `
	var resp {{.upperStartCamelObject}}
	err := m.conn.QueryRow(&resp, query, {{.lowerStartCamelPrimaryKey}})
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
func (m *{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}({{.in}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowIndex(&resp, {{.cacheKeyVariable}}, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := ` + "`" + `select ` + "`" + ` + {{.lowerStartCamelObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.originalField}} = ? limit 1` + "`" + `
		if err := conn.QueryRow(&resp, query, {{.lowerStartCamelField}}); err != nil {
			return nil, err
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := ` + "`" + `select ` + "`" + ` + {{.lowerStartCamelObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.originalPrimaryField}} = ? limit 1` + "`" + `
		return conn.QueryRow(v, query, primary)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{else}}var resp {{.upperStartCamelObject}}
	query := ` + "`" + `select ` + "`" + ` + {{.lowerStartCamelObject}}Rows + ` + "`" + ` from ` + "` + " + `m.table ` + " + `" + ` where {{.originalField}} limit 1` + "`" + `
	err := m.conn.QueryRow(&resp, query, {{.lowerStartCamelField}})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{end}}
`
