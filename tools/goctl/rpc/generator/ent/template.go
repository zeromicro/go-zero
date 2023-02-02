package ent

import (
	_ "embed"
)

var (
	//go:embed createOrUpdateLogic.tpl
	createOrUpdateTpl string

	//go:embed getListLogic.tpl
	getListLogicTpl string

	//go:embed deleteLogic.tpl
	deleteLogicTpl string

	//go:embed batchDeleteLogic.tpl
	batchDeleteLogicTpl string

	//go:embed updateStatusLogic.tpl
	updateStatusLogicTpl string

	//go:embed proto.tpl
	protoTpl string

	//go:embed pagination.tmpl
	PaginationTpl string
)
