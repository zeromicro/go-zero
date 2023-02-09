package ent

import (
	_ "embed"
)

var (
	//go:embed createLogic.tpl
	createTpl string

	//go:embed updateLogic.tpl
	updateTpl string

	//go:embed getListLogic.tpl
	getListLogicTpl string

	//go:embed getByIdLogic.tpl
	getByIdLogicTpl string

	//go:embed deleteLogic.tpl
	deleteLogicTpl string

	//go:embed proto.tpl
	protoTpl string

	//go:embed pagination.tmpl
	PaginationTpl string

	//go:embed notempty.tmpl
	NotEmptyTpl string
)
