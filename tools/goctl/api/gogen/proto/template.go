package proto

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

	//go:embed api.tpl
	apiTpl string
)
