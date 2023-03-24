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
)
