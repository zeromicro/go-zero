package template

import "fmt"

// GenFindOne defines find row by id.
func GenFindOne(dialect SqlDialect) string {
	p1 := dialect.PositionalParameter(0)
	return fmt.Sprintf(`
func (m *default{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRow(&resp, {{.cacheKeyVariable}}, func(conn sqlx.SqlConn, v interface{}) error {
		query :=  fmt.Sprintf("select %%s from %%s where {{.originalPrimaryKey}} = %s limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		return conn.QueryRow(v, query, {{.lowerStartCamelPrimaryKey}})
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}{{else}}query := fmt.Sprintf("select %%s from %%s where {{.originalPrimaryKey}} = %s limit 1", {{.lowerStartCamelObject}}Rows, m.table)
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
`, p1, p1)
}

// GenFindOneByField defines find row by field.
func GenFindOneByField(dialect SqlDialect) string {
	return fmt.Sprintf(`
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}({{.in}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowIndex(&resp, {{.cacheKeyVariable}}, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %%s from %%s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		if err := conn.QueryRow(&resp, query, {{.lowerStartCamelField}}); err != nil {
			return nil, err
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{else}}var resp {{.upperStartCamelObject}}
	query := fmt.Sprintf("select %%s from %%s where {{.originalField}} limit 1", {{.lowerStartCamelObject}}Rows, m.table )
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
`)
}

// GenFindOneByFieldExtraMethod defines find row by field with extras.
func GenFindOneByFieldExtraMethod(dialect SqlDialect) string {
	p1 := dialect.PositionalParameter(0)
	return fmt.Sprintf(`
func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%%s%%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %%s from %%s where {{.originalPrimaryField}} = %s limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	return conn.QueryRow(v, query, primary)
}
`, p1)
}

// FindOneMethod defines find row method.
var FindOneMethod = `FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error)`

// FindOneByFieldMethod defines find row by field method.
var FindOneByFieldMethod = `FindOneBy{{.upperField}}({{.in}}) (*{{.upperStartCamelObject}}, error) `
