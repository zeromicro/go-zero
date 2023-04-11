package template

import _ "embed"

var (
	//go:embed tmpl/not_empty_update.tmpl
	NotEmptyTmpl string

	//go:embed tmpl/pagination.tmpl
	PaginationTmpl string
)
