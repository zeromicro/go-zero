package template

// TableName defines a template that generate the tableName method.
const TableName = `
func (m *default{{.upperStartCamelObject}}Model) tableName() string {
	return m.table
}
`
