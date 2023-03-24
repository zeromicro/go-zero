package enttemplate

import _ "embed"

var (
	//go:embed pagination.tmpl
	PaginationTpl string

	//go:embed notempty.tmpl
	NotEmptyTpl string
)
